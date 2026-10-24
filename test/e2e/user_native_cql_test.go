package e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/elsyahtech/gorest-api-starter/app/handler"
	"github.com/elsyahtech/gorest-api-starter/app/library"
	"github.com/elsyahtech/gorest-api-starter/app/model"
	"github.com/elsyahtech/gorest-api-starter/app/router"
	"github.com/elsyahtech/gorest/database"
	gorestlog "github.com/elsyahtech/gorest/log"
	"github.com/elsyahtech/gorest/server"
)

func TestCreateUserWithNativeCQLE2E(t *testing.T) {
	testServer, app := newScyllaUserTestServer(t)
	expected := entity.User{
		ID:           "user-native-cql-e2e-01",
		Email:        "native-cql-e2e@example.test",
		Password:     "plain-text-native-cql-password",
		DepartmentID: "department-native-cql-e2e-01",
		IsActive:     true,
	}
	body, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("encode create user request: %v", err)
	}

	createResponse := sendRequest(
		t,
		testServer,
		http.MethodPost,
		"/examples/native-queries/users/scylla",
		body,
	)
	if createResponse.StatusCode < http.StatusOK || createResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("POST /examples/native-queries/users/scylla status = %d, want 2xx; body: %s",
			createResponse.StatusCode, createResponse.Body)
	}

	var createResult struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(createResponse.Body, &createResult); err != nil {
		t.Fatalf("decode native CQL create response: %v", err)
	}
	if !createResult.Success {
		t.Fatal("native CQL create response success = false, want true")
	}

	persisted := readScyllaUser(t, app, expected.ID)
	assertNativeCQLUserPersisted(t, persisted, expected)

	duplicate := expected
	duplicate.Email = "changed-native-cql-e2e@example.test"
	duplicate.Password = "another-native-cql-password"
	duplicate.DepartmentID = "department-native-cql-e2e-changed"
	duplicate.IsActive = false
	duplicateBody, err := json.Marshal(duplicate)
	if err != nil {
		t.Fatalf("encode duplicate user request: %v", err)
	}
	duplicateResponse := sendRequest(
		t,
		testServer,
		http.MethodPost,
		"/examples/native-queries/users/scylla",
		duplicateBody,
	)
	if duplicateResponse.StatusCode != http.StatusConflict {
		t.Fatalf("POST duplicate CQL user status = %d, want %d; body: %s",
			duplicateResponse.StatusCode, http.StatusConflict, duplicateResponse.Body)
	}

	persistedAfterDuplicate := readScyllaUser(t, app, expected.ID)
	if persistedAfterDuplicate != persisted {
		t.Fatalf("duplicate create changed persisted user: before=%+v after=%+v",
			persisted, persistedAfterDuplicate)
	}
}

func readScyllaUser(t *testing.T, app *gorest.App, userID string) entity.User {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	iter, _, _, err := app.Database().QueryCQL(ctx, `
		SELECT id, email, password, department_id, is_active, created_at
		FROM users WHERE id = ?
	`, userID)
	if err != nil {
		t.Fatalf("query persisted ScyllaDB user: %v", err)
	}

	var persisted entity.User
	if !iter.Scan(
		&persisted.ID,
		&persisted.Email,
		&persisted.Password,
		&persisted.DepartmentID,
		&persisted.IsActive,
		&persisted.CreatedAt,
	) {
		if err := iter.Close(); err != nil {
			t.Fatalf("close ScyllaDB user query: %v", err)
		}
		t.Fatal("ScyllaDB user query returned no row")
	}
	if err := iter.Close(); err != nil {
		t.Fatalf("close ScyllaDB user query: %v", err)
	}
	return persisted
}

func assertNativeCQLUserPersisted(t *testing.T, persisted, expected entity.User) {
	t.Helper()
	if persisted.ID != expected.ID ||
		persisted.Email != expected.Email ||
		persisted.DepartmentID != expected.DepartmentID ||
		persisted.IsActive != expected.IsActive {
		t.Fatalf("persisted ScyllaDB user = %+v, does not match request", persisted)
	}
	if persisted.Password == "" || persisted.Password == expected.Password {
		t.Fatal("persisted ScyllaDB password is empty or plaintext; want a hashed password")
	}
	if persisted.CreatedAt.IsZero() {
		t.Fatal("persisted ScyllaDB createdAt is zero, want a creation timestamp")
	}
}

func newScyllaUserTestServer(t *testing.T) (*server.Server, *gorest.App) {
	t.Helper()
	if os.Getenv("GOREST_SCYLLA_E2E") != "1" {
		t.Skip("set GOREST_SCYLLA_E2E=1 to run the ScyllaDB-backed E2E test")
	}

	host := os.Getenv("GOREST_SCYLLA_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("GOREST_SCYLLA_PORT")
	if port == "" {
		port = "9042"
	}
	username := driverTestSetting("GOREST_SCYLLA_USERNAME", "GOREST_DB_USERNAME")
	password := driverTestSetting("GOREST_SCYLLA_PASSWORD", "GOREST_DB_PASSWORD")
	dataDir := t.TempDir()
	logConfig := gorestlog.ConfigDefault
	logConfig.LogDir = filepath.Join(dataDir, "logs")
	keyspace := fmt.Sprintf("gorest_user_e2e_%x", time.Now().UnixNano())

	bootstrapApp := newScyllaTestApp(t, dataDir, logConfig, host, port, username, password, "system")
	bootstrapReady := false
	t.Cleanup(func() {
		if bootstrapReady {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if _, _, err := bootstrapApp.Database().ExecCQL(ctx, "DROP KEYSPACE IF EXISTS "+keyspace); err != nil {
				t.Errorf("drop ScyllaDB E2E keyspace: %v", err)
			}
		}
		if message, _, err := bootstrapApp.Close(); err != nil {
			t.Errorf("close ScyllaDB bootstrap application: %s: %v", message, err)
		}
	})
	if err := bootstrapApp.RunDatabaseStartupTasks(); err != nil {
		t.Fatalf("initialize ScyllaDB bootstrap application: %v", err)
	}
	bootstrapReady = true

	if err := executeScyllaCQL(t, bootstrapApp, fmt.Sprintf(
		"CREATE KEYSPACE %s WITH replication = {'class': 'NetworkTopologyStrategy', 'datacenter1': 1}", keyspace,
	)); err != nil {
		t.Fatalf("create ScyllaDB E2E keyspace: %v", err)
	}
	if err := executeScyllaCQL(t, bootstrapApp, fmt.Sprintf(`
		CREATE TABLE %s.users (
			id text PRIMARY KEY,
			email text,
			password text,
			department_id text,
			is_active boolean,
			created_at timestamp,
			updated_at timestamp,
			deleted_at timestamp
		)
	`, keyspace)); err != nil {
		t.Fatalf("create ScyllaDB E2E users table: %v", err)
	}

	app := newScyllaTestApp(t, dataDir, logConfig, host, port, username, password, keyspace)
	t.Cleanup(func() {
		if message, _, err := app.Close(); err != nil {
			t.Errorf("close ScyllaDB test application: %s: %v", message, err)
		}
	})
	if err := app.RunDatabaseStartupTasks(); err != nil {
		t.Fatalf("initialize ScyllaDB test application: %v", err)
	}

	testServer, _, err := server.New(&server.ConfigDefault)
	if err != nil {
		t.Fatalf("create test HTTP server: %v", err)
	}
	userModel := model.NewUserModel(app)
	userLibrary := library.NewUserLibrary(userModel)
	userHandler := handler.NewUserHandler(app, userLibrary)
	router.RegisterUserRoutes(testServer.Group("/examples"), userHandler)

	return testServer, app
}

func TestDriverTestSetting(t *testing.T) {
	t.Setenv("GOREST_DB_USERNAME", "app-user")
	t.Setenv("GOREST_DB_PASSWORD", "app-test-password")
	t.Setenv("GOREST_MONGO_USERNAME", "mongo-test-user")

	if got := driverTestSetting("GOREST_MONGO_USERNAME", "GOREST_DB_USERNAME"); got != "mongo-test-user" {
		t.Fatalf("driver-specific setting = %q, want override %q", got, "mongo-test-user")
	}
	if got := driverTestSetting("GOREST_SCYLLA_PASSWORD", "GOREST_DB_PASSWORD"); got != "app-test-password" {
		t.Fatalf("driver credential fallback = %q, want application database credential", got)
	}
}

func driverTestSetting(testKey, appKey string) string {
	if value, ok := os.LookupEnv(testKey); ok {
		return value
	}
	return os.Getenv(appKey)
}

func executeScyllaCQL(t *testing.T, app *gorest.App, query string, args ...any) error {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, _, err := app.Database().ExecCQL(ctx, query, args...)
	return err
}

func newScyllaTestApp(
	t *testing.T,
	dataDir string,
	logConfig gorestlog.Config,
	host, port, username, password, keyspace string,
) *gorest.App {
	t.Helper()
	databaseConfig := database.Config{
		Driver:            database.SCYLLA,
		Host:              host,
		Port:              port,
		Username:          username,
		Password:          password,
		Name:              keyspace,
		DatabaseDirectory: dataDir,
		Migration:         false,
		Seeder:            false,
		Disable:           false,
		Timeout:           10 * time.Second,
	}
	app := gorest.New(gorest.Config{
		AppName:     "gorest-api-starter-user-scylla-e2e",
		Version:     "test",
		Environment: "Test",
		Timezone:    "UTC",
	})
	app.Use(
		gorestlog.Run(logConfig),
		database.Run(databaseConfig),
	)
	return app
}
