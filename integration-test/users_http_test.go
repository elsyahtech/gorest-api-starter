//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
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
)

func TestPostgreSQLUserEndpointsOverHTTP(t *testing.T) {
	testConfig, err := loadPostgreSQLTestConfig()
	if err != nil {
		t.Fatal(err)
	}
	applyPostgreSQLTestConfig(t, testConfig)

	dbConfig, err := config.DatabaseConfig()
	if err != nil {
		t.Fatalf("load test database configuration: %v", err)
	}
	if dbConfig.Driver != database.POSTGRES || dbConfig.Name != testConfig.database {
		t.Fatalf("database configuration is not the selected PostgreSQL test database")
	}
	databaseContext, cancelDatabase := context.WithTimeout(context.Background(), 15*time.Second)
	testDB, message, err := database.New(databaseContext, dbConfig, gorest.DefaultTimezone)
	cancelDatabase()
	if err != nil {
		t.Fatalf("connect to configured PostgreSQL test database (%s): %v", message, err)
	}
	t.Cleanup(func() {
		if err := testDB.SQL.Close(); err != nil {
			t.Errorf("close PostgreSQL test database connection: %v", err)
		}
	})
	// Release the preflight connection before the application uses the test account.
	testDB.SQL.SetMaxIdleConns(0)

	signatureKeyBytes := make([]byte, 48)
	if _, err := rand.Read(signatureKeyBytes); err != nil {
		t.Fatalf("generate temporary token signing key: %v", err)
	}
	signatureKey := base64.RawURLEncoding.EncodeToString(signatureKeyBytes)
	baseURL, stopApp := startApplication(t, testConfig, signatureKey)
	defer stopApp()
	requireUsersTable(t, testDB.SQL)
	requirePostgreSQLSeedData(t, testDB.SQL)

	testID := uuid.NewString()
	testEmail := "integration+" + strings.ReplaceAll(testID, "-", "") + "@example.test"
	testPassword := "integration-test-password"
	cleanupUser(t, testDB.SQL, testID, testEmail)
	t.Cleanup(func() {
		cleanupUser(t, testDB.SQL, testID, testEmail)
	})

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
	t.Run("creates a PostgreSQL user", func(t *testing.T) {
		response := requestJSON(t, client, http.MethodPost, baseURL+"/users/create", createRequest, "")
		if response.status != http.StatusOK || response.json["success"] != true {
			t.Fatalf("create response = (%d, %s), want HTTP 200 success", response.status, response.body)
		}

		var storedID, storedEmail, storedPassword string
		var active bool
		err := testDB.SQL.QueryRow(
			"SELECT id, email, password, is_active FROM users WHERE id = $1",
			testID,
		).Scan(&storedID, &storedEmail, &storedPassword, &active)
		if err != nil {
			t.Fatalf("read created user from PostgreSQL: %v", err)
		}
		if storedID != testID || storedEmail != testEmail || !active || storedPassword == testPassword {
			t.Fatalf("persisted user = (%q, %q, %t, password-set:%t), want created user with hashed password",
				storedID, storedEmail, active, storedPassword != "")
		}
	})

	t.Run("rejects duplicate user records", func(t *testing.T) {
		response := requestJSON(t, client, http.MethodPost, baseURL+"/users/create", createRequest, "")
		if response.status < http.StatusBadRequest || response.json["success"] != false {
			t.Fatalf("duplicate create response = (%d, %s), want a client/server error", response.status, response.body)
		}

		var count int
		if err := testDB.SQL.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", testID).Scan(&count); err != nil {
			t.Fatalf("count duplicate test user rows: %v", err)
		}
		if count != 1 {
			t.Fatalf("duplicate create left %d records for test ID, want exactly one", count)
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

func requireUsersTable(t *testing.T, db interface {
	QueryRow(query string, args ...any) *sql.Row
}) {
	t.Helper()

	var tableCount int
	if err := db.QueryRow(
		"SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = current_schema() AND table_name = 'users'",
	).Scan(&tableCount); err != nil {
		t.Fatalf("check the PostgreSQL users table: %v", err)
	}
	if tableCount != 1 {
		t.Fatal("PostgreSQL test database does not contain the users table after running the migration")
	}
}

func requirePostgreSQLSeedData(t *testing.T, db interface {
	QueryRow(query string, args ...any) *sql.Row
}) {
	t.Helper()

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM users WHERE id LIKE '3746e32f-1bd6-4023-a608-f535696f5e%'").Scan(&count); err != nil {
		t.Fatalf("verify PostgreSQL seed data: %v", err)
	}
	if count != 10 {
		t.Fatalf("PostgreSQL seeder created %d sample users, want 10", count)
	}
}

func cleanupUser(t *testing.T, db interface {
	Exec(query string, args ...any) (sql.Result, error)
}, id, email string) {
	t.Helper()

	if _, err := db.Exec("DELETE FROM users WHERE id = $1 AND email = $2", id, email); err != nil {
		t.Errorf("clean up integration user %s: %v", id, err)
	}
}

func applyPostgreSQLTestConfig(t *testing.T, config postgresTestConfig) {
	t.Helper()
	for key, value := range map[string]string{
		"GOREST_APP_ENV":                      "testing",
		"GOREST_DB_DISABLE":                   "false",
		"GOREST_DB_DRIVER":                    database.POSTGRES,
		"GOREST_DB_HOST":                      config.host,
		"GOREST_DB_PORT":                      config.port,
		"GOREST_DB_NAME":                      config.database,
		"GOREST_DB_USER":                      config.username,
		"GOREST_DB_PASSWORD":                  config.password,
		"GOREST_DB_TLS":                       config.tls,
		"GOREST_DB_CA_PATH":                   config.caPath,
		"GOREST_DB_CERT_PATH":                 config.certPath,
		"GOREST_DB_KEY_PATH":                  config.keyPath,
		"GOREST_DB_DIRECTORY":                 config.databaseDirectory,
		"GOREST_DB_MIGRATION":                 "true",
		"GOREST_DB_SEEDER":                    "true",
		"GOREST_DB_ALLOW_INSECURE_PRODUCTION": "false",
	} {
		t.Setenv(key, value)
	}
}

func startApplication(t *testing.T, config postgresTestConfig, signatureKey string) (string, func()) {
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
	deadline := time.Now().Add(30 * time.Second)
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
