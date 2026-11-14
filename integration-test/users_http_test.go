//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/config"
	"github.com/elsyahtech/gorest/database"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestMongoUserEndpointsOverHTTP(t *testing.T) {
	testConfig, err := loadMongoTestConfig()
	if err != nil {
		t.Fatal(err)
	}
	prepareMongoIntegrationSchema(t, testConfig)

	dbConfig, err := config.DatabaseConfig()
	if err != nil {
		t.Fatalf("load test database configuration: %v", err)
	}
	if dbConfig.Driver != database.MONGO || dbConfig.Name != testConfig.database {
		t.Fatal("database configuration is not the selected MongoDB test database")
	}
	databaseContext, cancelDatabase := context.WithTimeout(context.Background(), 20*time.Second)
	testDB, message, err := database.New(databaseContext, dbConfig, gorest.DefaultTimezone)
	cancelDatabase()
	if err != nil {
		t.Fatalf("connect to configured MongoDB test database (%s): %v", message, err)
	}
	t.Cleanup(func() {
		cleanupMongoTestConnection(t, testDB.Mongo)
	})
	collection := testDB.Mongo.Database(testConfig.database).Collection("users")
	requireUsersCollection(t, testDB.Mongo, testConfig.database)

	testID := uuid.NewString()
	testEmail := "integration+" + strings.ReplaceAll(testID, "-", "") + "@example.test"
	testPassword := "integration-test-password"
	cleanupMongoUser(t, collection, testID, testEmail)
	t.Cleanup(func() {
		cleanupMongoUser(t, collection, testID, testEmail)
	})

	signatureKeyBytes := make([]byte, 48)
	if _, err := rand.Read(signatureKeyBytes); err != nil {
		t.Fatalf("generate temporary token signing key: %v", err)
	}
	signatureKey := base64.RawURLEncoding.EncodeToString(signatureKeyBytes)
	baseURL, stopApp := startApplication(t, testConfig, signatureKey)
	defer stopApp()

	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("rejects malformed create request", func(t *testing.T) {
		responseBody := requestJSON(t, client, http.MethodPost, baseURL+"/users/create", "{", "")
		if responseBody.status != http.StatusBadRequest || responseBody.json["success"] != false {
			t.Fatalf("malformed request response = (%d, %s), want a 400 error envelope", responseBody.status, responseBody.body)
		}
	})
	t.Run("rejects missing required user fields", func(t *testing.T) {
		response := requestJSON(t, client, http.MethodPost, baseURL+"/users/create", `{}`, "")
		if response.status != http.StatusBadRequest || response.json["success"] != false {
			t.Fatalf("missing-field response = (%d, %s), want HTTP 400 error", response.status, response.body)
		}
	})

	createRequest := fmt.Sprintf(
		`{"id":%q,"email":%q,"password":%q,"isActive":true}`,
		testID,
		testEmail,
		testPassword,
	)
	t.Run("creates a MongoDB user", func(t *testing.T) {
		response := requestJSON(t, client, http.MethodPost, baseURL+"/users/create", createRequest, "")
		if response.status != http.StatusOK || response.json["success"] != true {
			t.Fatalf("create response = (%d, %s), want HTTP 200 success", response.status, response.body)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var stored struct {
			ID       string `bson:"id"`
			Email    string `bson:"email"`
			Password string `bson:"password"`
			IsActive bool   `bson:"is_active"`
		}
		if err := collection.FindOne(ctx, bson.M{"id": testID}).Decode(&stored); err != nil {
			t.Fatalf("read created user from MongoDB: %v", err)
		}
		if stored.ID != testID || stored.Email != testEmail || !stored.IsActive || stored.Password == testPassword {
			t.Fatalf("persisted MongoDB user = (%q, %q, %t, password-set:%t), want created user with hashed password",
				stored.ID, stored.Email, stored.IsActive, stored.Password != "")
		}
	})

	t.Run("rejects duplicate user records", func(t *testing.T) {
		response := requestJSON(t, client, http.MethodPost, baseURL+"/users/create", createRequest, "")
		if response.status < http.StatusBadRequest || response.json["success"] != false {
			t.Fatalf("duplicate create response = (%d, %s), want a client/server error", response.status, response.body)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		count, err := collection.CountDocuments(ctx, bson.M{"id": testID})
		if err != nil {
			t.Fatalf("count duplicate test user documents: %v", err)
		}
		if count != 1 {
			t.Fatalf("duplicate create left %d documents for test ID, want exactly one", count)
		}
	})

	t.Run("rejects incorrect password", func(t *testing.T) {
		request := fmt.Sprintf(`{"email":%q,"password":%q}`, testEmail, "wrong-password")
		response := requestJSON(t, client, http.MethodPost, baseURL+"/users/login", request, "")
		if response.status != http.StatusUnauthorized || response.json["success"] != false {
			t.Fatalf("wrong-password response = (%d, %s), want HTTP 401 error", response.status, response.body)
		}
	})

	var accessCookie *http.Cookie
	t.Run("logs in and issues an access cookie", func(t *testing.T) {
		request := fmt.Sprintf(`{"email":%q,"password":%q}`, testEmail, testPassword)
		response := requestJSON(t, client, http.MethodPost, baseURL+"/users/login", request, "")
		if response.status != http.StatusOK || response.json["success"] != true {
			t.Fatalf("login response = (%d, %s), want HTTP 200 success", response.status, response.body)
		}
		for _, cookie := range response.cookies {
			if cookie.Name == "access" {
				accessCookie = cookie
				break
			}
		}
		if accessCookie == nil || accessCookie.Value == "" {
			t.Fatalf("login response did not issue an access cookie: %v", response.cookies)
		}
	})

	t.Run("requires authentication to list users", func(t *testing.T) {
		response := requestJSON(t, client, http.MethodGet, baseURL+"/users/", "", "")
		if response.status != http.StatusUnauthorized || response.json["Error"] == "" {
			t.Fatalf("unauthenticated list response = (%d, %s), want HTTP 401 error", response.status, response.body)
		}
	})

	t.Run("lists users without exposing password hashes", func(t *testing.T) {
		response := requestJSON(t, client, http.MethodGet, baseURL+"/users/", "", accessCookie.String())
		if response.status != http.StatusOK || response.json["success"] != true {
			t.Fatalf("authenticated list response = (%d, %s), want HTTP 200 success", response.status, response.body)
		}
		if !strings.Contains(response.body, testEmail) || strings.Contains(response.body, "password") {
			t.Fatalf("user-list response does not safely contain the created user: %s", response.body)
		}
	})
}

type httpResponse struct {
	status  int
	body    string
	json    map[string]any
	cookies []*http.Cookie
}

func requestJSON(t *testing.T, client *http.Client, method, url, body, cookie string) httpResponse {
	t.Helper()

	request, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		t.Fatalf("create %s request: %v", method, err)
	}
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	if cookie != "" {
		request.Header.Set("Cookie", cookie)
	}

	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("%s %s request failed: %v", method, url, err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		t.Fatalf("read %s response: %v", method, err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(responseBody, &decoded); err != nil {
		t.Fatalf("decode %s response JSON %q: %v", method, responseBody, err)
	}

	return httpResponse{
		status:  response.StatusCode,
		body:    string(responseBody),
		json:    decoded,
		cookies: response.Cookies(),
	}
}

func requireUsersCollection(t *testing.T, client *mongo.Client, databaseName string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collections, err := client.Database(databaseName).ListCollectionNames(ctx, bson.D{{Key: "name", Value: "users"}})
	if err != nil {
		t.Fatalf("check the MongoDB users collection: %v", err)
	}
	if len(collections) != 1 || collections[0] != "users" {
		t.Fatal("MongoDB test database does not contain the users collection; apply the MongoDB migration before running integration tests")
	}
}

func applyMongoTestConfig(t *testing.T, config mongoTestConfig) {
	t.Helper()
	for key, value := range map[string]string{
		"GOREST_APP_ENV":                      "testing",
		"GOREST_DB_DISABLE":                   "false",
		"GOREST_DB_DRIVER":                    database.MONGO,
		"GOREST_DB_HOST":                      config.host,
		"GOREST_DB_PORT":                      config.port,
		"GOREST_DB_NAME":                      config.database,
		"GOREST_DB_USER":                      config.username,
		"GOREST_DB_PASSWORD":                  config.password,
		"GOREST_DB_AUTH_SOURCE":               config.authSource,
		"GOREST_DB_TLS":                       config.tls,
		"GOREST_DB_CA_PATH":                   config.caPath,
		"GOREST_DB_CERT_PATH":                 config.certPath,
		"GOREST_DB_KEY_PATH":                  config.keyPath,
		"GOREST_DB_DIRECTORY":                 config.databaseDir,
		"GOREST_DB_MIGRATION":                 "false",
		"GOREST_DB_SEEDER":                    "false",
		"GOREST_DB_ALLOW_INSECURE_PRODUCTION": "false",
	} {
		t.Setenv(key, value)
	}
}

func prepareMongoIntegrationSchema(t *testing.T, testConfig mongoTestConfig) {
	t.Helper()
	applyMongoTestConfig(t, testConfig)

	dbConfig, err := config.DatabaseConfig()
	if err != nil {
		t.Fatalf("load MongoDB migration configuration: %v", err)
	}
	dbConfig.Migration = true
	dbConfig.Seeder = true

	app := gorest.New(gorest.ConfigDefault)
	app.Use(database.Run(*dbConfig))
	t.Cleanup(func() {
		if _, _, err := app.Close(); err != nil {
			t.Errorf("close MongoDB migration application: %v", err)
		}
	})
	if err := app.RunDatabaseStartupTasks(); err != nil {
		t.Fatalf("apply MongoDB integration migration: %v", err)
	}
}

func startApplication(t *testing.T, config mongoTestConfig, signatureKey string) (string, func()) {
	t.Helper()

	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve an HTTP port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		t.Fatalf("release reserved HTTP port: %v", err)
	}

	projectRoot, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("resolve project root: %v", err)
	}
	temporaryDirectory := t.TempDir()
	binaryPath := filepath.Join(temporaryDirectory, "gorest-api-starter")
	buildArgs := []string{"build"}
	coverageDirectory := strings.TrimSpace(os.Getenv("GOREST_INTEGRATION_COVERAGE_DIR"))
	if coverageDirectory != "" {
		buildArgs = append(buildArgs, "-cover", "-coverpkg=.,./app/...")
	}
	buildArgs = append(buildArgs, "-o", binaryPath, ".")
	build := exec.Command("go", buildArgs...)
	build.Dir = projectRoot
	var buildOutput bytes.Buffer
	build.Stdout = &buildOutput
	build.Stderr = &buildOutput
	if err := build.Run(); err != nil {
		t.Fatalf("build application for integration test: %v\n%s", err, buildOutput.String())
	}

	logDirectory := filepath.Join(temporaryDirectory, "logs")
	if err := os.MkdirAll(logDirectory, 0o700); err != nil {
		t.Fatalf("create temporary application log directory: %v", err)
	}
	environment := config.appEnvironment(fmt.Sprint(port), logDirectory, signatureKey)
	if coverageDirectory != "" {
		environment["GOCOVERDIR"] = coverageDirectory
	}
	command := exec.Command(binaryPath)
	command.Dir = temporaryDirectory
	command.Env = mergeEnvironment(os.Environ(), environment)
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Start(); err != nil {
		t.Fatalf("start application: %v", err)
	}

	processDone := make(chan struct{})
	var processErr error
	go func() {
		processErr = command.Wait()
		close(processDone)
	}()

	var stopOnce sync.Once
	stop := func() {
		stopOnce.Do(func() {
			select {
			case <-processDone:
				return
			default:
			}
			if err := command.Process.Signal(os.Interrupt); err != nil {
				t.Errorf("send graceful shutdown signal to application: %v", err)
				_ = command.Process.Kill()
			}
			select {
			case <-processDone:
			case <-time.After(10 * time.Second):
				_ = command.Process.Kill()
				<-processDone
				t.Error("application did not stop after graceful shutdown")
			}
		})
	}
	t.Cleanup(stop)

	baseURL := fmt.Sprintf("http://127.0.0.1:%d", port)
	client := &http.Client{Timeout: time.Second}
	deadline := time.Now().Add(45 * time.Second)
	ready := false
	for time.Now().Before(deadline) {
		select {
		case <-processDone:
			t.Fatalf("application exited before becoming ready: %v\nstdout:\n%s\nstderr:\n%s",
				processErr, stdout.String(), stderr.String())
		default:
		}

		response, requestErr := client.Get(baseURL + "/gorest/ready")
		if requestErr == nil {
			_ = response.Body.Close()
			if response.StatusCode == http.StatusOK {
				ready = true
				break
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
	if !ready {
		t.Fatalf("application did not become ready at %s\nstdout:\n%s\nstderr:\n%s",
			baseURL, stdout.String(), stderr.String())
	}

	return baseURL, stop
}

func mergeEnvironment(environment []string, overrides map[string]string) []string {
	values := make(map[string]string, len(environment)+len(overrides))
	for _, entry := range environment {
		key, value, ok := strings.Cut(entry, "=")
		if ok && !strings.HasPrefix(key, "GOREST_") {
			values[key] = value
		}
	}
	for key, value := range overrides {
		values[key] = value
	}
	merged := make([]string, 0, len(values))
	for key, value := range values {
		merged = append(merged, key+"="+value)
	}
	sort.Strings(merged)
	return merged
}
