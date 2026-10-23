package drivers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestDatabaseDriverIntegration(t *testing.T) {
	driver := strings.ToLower(strings.TrimSpace(os.Getenv("GOREST_DRIVER_INTEGRATION")))
	if driver == "" {
		t.Skip("set GOREST_DRIVER_INTEGRATION to run a disposable database integration test")
	}

	config, databaseName := integrationDatabaseConfig(t, driver)
	if driver == database.SQLSERVER {
		ensureSQLServerDatabase(t, config, databaseName)
		config.Name = databaseName
	}
	if driver == database.SCYLLA {
		ensureScyllaKeyspace(t, config, databaseName)
		config.Name = databaseName
	}

	app := newIntegrationApp(t, config)
	if err := app.RunDatabaseStartupTasks(); err != nil {
		t.Fatalf("run %s migrations and seeders: %v", driver, err)
	}

	verifySeedData(t, app)
	verifyCRUD(t, app, driver)
	if driver != database.MONGO && driver != database.SCYLLA {
		verifyTransactionRollback(t, app, driver)
		verifyMidStepRollback(t, app, driver)
	} else {
		t.Logf("%s transaction atomicity: not applicable to the starter's GOREST ORM transaction API", driver)
	}
	verifyConflictAndInvalidRequests(t, app, driver)
	verifyConcurrentRequests(t, app, driver)
	verifyMigrationStartupFailure(t, config, driver)
}

func integrationDatabaseConfig(t *testing.T, driver string) (database.Config, string) {
	t.Helper()

	repositoryDatabaseDir, err := filepath.Abs(filepath.Join("..", "..", "app", "database"))
	if err != nil {
		t.Fatalf("resolve database example directory: %v", err)
	}

	config := database.Config{
		Driver:                 driver,
		Host:                   envOr("GOREST_DRIVER_HOST", "127.0.0.1"),
		DatabaseDirectory:      repositoryDatabaseDir,
		MigrationDirectory:     "migrations",
		SeederDirectory:        "seeders",
		Timeout:                30 * time.Second,
		ConnectionMaxLifetime:  time.Minute,
		MaxOpenConnections:     5,
		MaxIdleConnections:     2,
		Enabled:                true,
		Migration:              true,
		Seeder:                 true,
		TrustServerCertificate: true,
	}
	databaseName := envOr("GOREST_DRIVER_DATABASE", "gorest_driver_it")

	switch driver {
	case database.SQLITE:
		// SQLite has a single-writer model; serialize DB work in the test pool
		// while keeping all HTTP requests concurrent.
		config.MaxOpenConnections = 1
		config.MaxIdleConnections = 1
		config.Name = filepath.Join(t.TempDir(), "gorest-driver-integration.sqlite")
		config.DatabaseDirectory = repositoryDatabaseDir
	case database.POSTGRES:
		config.Port = envOr("GOREST_DRIVER_PORT", "15432")
		config.Username = envOr("GOREST_DRIVER_USERNAME", "gorest_test")
		config.Password = envOr("GOREST_DRIVER_PASSWORD", "GorestTest!Pass2026")
		config.Name = databaseName
		config.TLS = "disable"
	case database.MYSQL:
		config.Port = envOr("GOREST_DRIVER_PORT", "13306")
		config.Username = envOr("GOREST_DRIVER_USERNAME", "root")
		config.Password = envOr("GOREST_DRIVER_PASSWORD", "GorestTest!Pass2026")
		config.Name = databaseName
		config.TLS = "disable"
	case database.ORACLE:
		config.Host = envOr("GOREST_DRIVER_HOST", "127.0.0.1")
		config.Port = envOr("GOREST_DRIVER_PORT", "1521")
		config.Username = os.Getenv("GOREST_DRIVER_USERNAME")
		config.Password = os.Getenv("GOREST_DRIVER_PASSWORD")
		config.Name = envOr("GOREST_DRIVER_DATABASE", databaseName)
		config.LibDir = os.Getenv("GOREST_DRIVER_LIB_DIR")
		config.TLS = "disable"
	case database.SQLSERVER:
		config.Port = envOr("GOREST_DRIVER_PORT", "11433")
		config.Username = envOr("GOREST_DRIVER_USERNAME", "sa")
		config.Password = envOr("GOREST_DRIVER_PASSWORD", "GorestTest!Pass2026")
		config.Name = "master"
		config.TLS = "disable"
	case database.MONGO:
		config.Port = envOr("GOREST_DRIVER_PORT", "27018")
		config.Username = os.Getenv("GOREST_DRIVER_USERNAME")
		config.Password = os.Getenv("GOREST_DRIVER_PASSWORD")
		config.Name = databaseName
		config.TLS = "false"
	case database.SCYLLA:
		config.Port = envOr("GOREST_DRIVER_PORT", "19042")
		config.Username = os.Getenv("GOREST_DRIVER_USERNAME")
		config.Password = os.Getenv("GOREST_DRIVER_PASSWORD")
		config.Name = "system"
	default:
		t.Fatalf("unsupported integration driver %q", driver)
	}

	return config, databaseName
}

func newIntegrationApp(t *testing.T, config database.Config) *gorest.App {
	t.Helper()

	logConfig := gorestlog.ConfigDefault
	logConfig.LogDir = filepath.Join(t.TempDir(), "logs")
	app := gorest.New(gorest.Config{
		AppName:     "gorest-api-starter-driver-integration",
		Version:     "test",
		Environment: "Test",
		Timezone:    "UTC",
	})
	app.Use(
		gorestlog.Run(logConfig),
		database.Run(config),
	)
	t.Cleanup(func() {
		if message, _, err := app.Close(); err != nil {
			t.Errorf("close %s integration app (%s): %v", config.Driver, message, err)
		} else {
			t.Logf("%s application shutdown passed", config.Driver)
		}
	})
	return app
}

func ensureSQLServerDatabase(t *testing.T, config database.Config, databaseName string) {
	t.Helper()

	config.Migration = false
	config.Seeder = false
	app := newIntegrationApp(t, config)
	query := fmt.Sprintf(
		"IF DB_ID(N'%s') IS NULL CREATE DATABASE [%s]",
		strings.ReplaceAll(databaseName, "'", "''"),
		strings.ReplaceAll(databaseName, "]", "]]"),
	)
	ctx, cancel := app.NewContext(config.Timeout)
	defer cancel()
	if _, _, _, err := app.Database().ExecSQL(ctx, query); err != nil {
		t.Fatalf("create disposable SQL Server database: %v", err)
	}
}

func ensureScyllaKeyspace(t *testing.T, config database.Config, keyspace string) {
	t.Helper()

	config.Migration = false
	config.Seeder = false
	app := newIntegrationApp(t, config)
	ctx, cancel := app.NewContext(config.Timeout)
	defer cancel()
	query := fmt.Sprintf(
		"CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'NetworkTopologyStrategy', 'datacenter1': 1}",
		keyspace,
	)
	if _, _, err := app.Database().ExecCQL(ctx, query); err != nil {
		t.Fatalf("create disposable ScyllaDB keyspace: %v", err)
	}
}

func verifySeedData(t *testing.T, app *gorest.App) {
	t.Helper()

	var users []entity.User
	query := app.Database().Table("users").Find(&users)
	if query.Error != nil {
		t.Fatalf("read users created by the configured seeder: %s: %v", query.Message, query.Error)
	}
	if len(users) < 10 {
		t.Fatalf("seeded user count = %d, want at least 10", len(users))
	}
	t.Logf("migration and seeder passed; read %d seeded users", len(users))
}

func verifyCRUD(t *testing.T, app *gorest.App, driver string) {
	t.Helper()

	id := "driver-it-" + driver
	user := entity.User{
		CreatedAt:    time.Now().UTC(),
		ID:           id,
		Email:        id + "@example.test",
		Password:     "integration-test-hash",
		DepartmentID: "dept-01",
		IsActive:     true,
	}
	query := app.Database().Table("users").Create(&user)
	if query.Error != nil {
		t.Fatalf("create %s user through GOREST ORM: %s: %v", driver, query.Message, query.Error)
	}

	var created entity.User
	query = app.Database().Table("users").Where("id = ?", id).Find(&created)
	if query.Error != nil {
		t.Fatalf("read created %s user: %s: %v", driver, query.Message, query.Error)
	}
	if created.ID != id || created.Email != user.Email {
		t.Fatalf("read created %s user = %+v, want ID %q and email %q", driver, created, id, user.Email)
	}

	created.Email = "updated-" + id + "@example.test"
	query = app.Database().Table("users").Where("id = ?", id).Update(&created)
	if query.Error != nil {
		t.Fatalf("update %s user through GOREST ORM: %s: %v", driver, query.Message, query.Error)
	}

	var updated entity.User
	query = app.Database().Table("users").Where("id = ?", id).Find(&updated)
	if query.Error != nil {
		t.Fatalf("read updated %s user: %s: %v", driver, query.Message, query.Error)
	}
	if updated.Email != created.Email {
		t.Fatalf("updated %s email = %q, want %q", driver, updated.Email, created.Email)
	}

	query = app.Database().Table("users").Where("id = ?", id).Delete(&entity.User{})
	if query.Error != nil {
		t.Fatalf("delete %s user through GOREST ORM: %s: %v", driver, query.Message, query.Error)
	}

	var deleted entity.User
	query = app.Database().Table("users").Where("id = ?", id).Find(&deleted)
	if query.Error == nil {
		t.Fatalf("deleted %s user is still readable: %+v", driver, deleted)
	}
	t.Logf("%s ORM create/read/update/delete passed", driver)
}

func verifyTransactionRollback(t *testing.T, app *gorest.App, driver string) {
	t.Helper()

	id := "tx-it-" + driver
	tx, message, status, err := app.Database().BeginTx()
	if err != nil {
		t.Fatalf("begin %s transaction (HTTP %d, %s): %v", driver, status, message, err)
	}
	user := entity.User{
		CreatedAt:    time.Now().UTC(),
		ID:           id,
		Email:        id + "@example.test",
		Password:     "integration-test-hash",
		DepartmentID: "dept-01",
		IsActive:     true,
	}
	query := tx.Table("users").Create(&user)
	if query.Error != nil {
		_, _, _ = tx.Rollback()
		t.Fatalf("create user in %s transaction: %s: %v", driver, query.Message, query.Error)
	}
	if message, status, err := tx.Rollback(); err != nil {
		t.Fatalf("rollback %s transaction (HTTP %d, %s): %v", driver, status, message, err)
	}

	var rolledBack entity.User
	query = app.Database().Table("users").Where("id = ?", id).Find(&rolledBack)
	if query.Error == nil {
		t.Fatalf("%s rollback left transaction user persisted: %+v", driver, rolledBack)
	}
	t.Logf("%s transaction rollback passed", driver)
}

func verifyMidStepRollback(t *testing.T, app *gorest.App, driver string) {
	t.Helper()

	id := "tx-mid-step-" + driver
	duplicateEmail := id + "@example.test"
	tx, message, status, err := app.Database().BeginTx()
	if err != nil {
		t.Fatalf("begin %s mid-step transaction (HTTP %d, %s): %v", driver, status, message, err)
	}
	first := entity.User{
		CreatedAt: time.Now().UTC(), ID: id, Email: duplicateEmail,
		Password: "integration-test-hash", DepartmentID: "dept-01", IsActive: true,
	}
	if query := tx.Table("users").Create(&first); query.Error != nil {
		_, _, _ = tx.Rollback()
		t.Fatalf("create first %s transaction row: %s: %v", driver, query.Message, query.Error)
	}
	second := first
	second.ID = id + "-conflict"
	if query := tx.Table("users").Create(&second); query.Error == nil {
		_, _, _ = tx.Rollback()
		t.Fatalf("second %s transaction write with duplicate email succeeded; expected constraint failure", driver)
	}
	if message, status, err := tx.Rollback(); err != nil {
		t.Fatalf("rollback %s transaction after middle-step failure (HTTP %d, %s): %v", driver, status, message, err)
	}
	var rolledBack entity.User
	if query := app.Database().Table("users").Where("id = ?", id).Find(&rolledBack); query.Error == nil {
		t.Fatalf("%s mid-step rollback left the first transaction row persisted: %+v", driver, rolledBack)
	}
	t.Logf("%s mid-step constraint failure rolled back the earlier write", driver)
}

func verifyConflictAndInvalidRequests(t *testing.T, app *gorest.App, driver string) {
	t.Helper()
	server := newUserIntegrationServer(t, app)
	conflictID := fmt.Sprintf("dup-%x", time.Now().UnixNano())
	original := entity.User{CreatedAt: time.Now().UTC(), ID: conflictID, Email: conflictID + "@example.test", Password: "original", DepartmentID: "dept-01", IsActive: true}
	if query := app.Database().Table("users").Create(&original); query.Error != nil {
		t.Fatalf("create duplicate-key baseline for %s: %s: %v", driver, query.Message, query.Error)
	}
	duplicateKey := original
	duplicateKey.Email = "changed-" + original.Email
	query := app.Database().Table("users").Create(&duplicateKey)
	if query.Error == nil {
		t.Fatalf("%s duplicate primary key was accepted", driver)
	}
	var unchanged entity.User
	if read := app.Database().Table("users").Where("id = ?", conflictID).Find(&unchanged); read.Error != nil || unchanged.Email != original.Email {
		t.Fatalf("%s duplicate-key attempt changed baseline: %+v, query error = %v", driver, unchanged, read.Error)
	}
	if cleanup := app.Database().Table("users").Where("id = ?", conflictID).Delete(&entity.User{}); cleanup.Error != nil {
		t.Errorf("delete duplicate-key test row: %s: %v", cleanup.Message, cleanup.Error)
	}

	duplicateEmail := entity.UserArray{
		{ID: "driver-conflict-a", Email: "driver-request-conflict@example.test", Password: "test-password", DepartmentID: "dept-01", IsActive: true},
		{ID: "driver-conflict-b", Email: "driver-request-conflict@example.test", Password: "test-password", DepartmentID: "dept-01", IsActive: true},
	}
	body, err := json.Marshal(duplicateEmail)
	if err != nil {
		t.Fatalf("encode duplicate-email request: %v", err)
	}
	response := integrationRequest(t, server, http.MethodPost, "/examples/orm-queries/users/bulk", body)
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("duplicate-email request status = %d, want %d; body: %s", response.StatusCode, http.StatusBadRequest, response.Body)
	}
	assertDriverUserAbsent(t, app, "driver-conflict-a")
	assertDriverUserAbsent(t, app, "driver-conflict-b")

	response = integrationRequest(t, server, http.MethodPost, "/examples/orm-queries/users", []byte("{"))
	if response.StatusCode < http.StatusBadRequest || response.StatusCode >= http.StatusInternalServerError {
		t.Fatalf("malformed JSON request status = %d, want 4xx; body: %s", response.StatusCode, response.Body)
	}

	missingUpdate := entity.UserArray{{ID: "driver-missing-update", Email: "missing@example.test", Password: "test-password", DepartmentID: "dept-01", IsActive: true}}
	body, err = json.Marshal(missingUpdate)
	if err != nil {
		t.Fatalf("encode missing-user update: %v", err)
	}
	response = integrationRequest(t, server, http.MethodPut, "/examples/orm-queries/users/bulk", body)
	if response.StatusCode < http.StatusBadRequest {
		t.Fatalf("missing-user update status = %d, want an error response; body: %s", response.StatusCode, response.Body)
	}
	assertDriverUserAbsent(t, app, "driver-missing-update")
	t.Log("duplicate request key, malformed JSON, and missing-user update were rejected without creating rows")
}

func verifyConcurrentRequests(t *testing.T, app *gorest.App, driver string) {
	t.Helper()
	const requestCount = 8
	server := newUserIntegrationServer(t, app)
	start := make(chan struct{})
	type requestResult struct {
		id     string
		status int
		body   []byte
		err    error
	}
	results := make(chan requestResult, requestCount)
	var wait sync.WaitGroup
	runID := time.Now().UnixNano()
	ids := make([]string, requestCount)
	for i := 0; i < requestCount; i++ {
		id := fmt.Sprintf("c-%s-%x-%02d", driver, runID, i)
		ids[i] = id
		user := entity.User{ID: id, Email: id + "@example.test", Password: "test-password", DepartmentID: "dept-01", IsActive: true}
		body, err := json.Marshal(user)
		if err != nil {
			t.Fatalf("encode concurrent user request: %v", err)
		}
		wait.Add(1)
		go func(body []byte, id string) {
			defer wait.Done()
			<-start
			response, err := integrationRequestResult(server, http.MethodPost, "/examples/orm-queries/users", body)
			if err != nil {
				results <- requestResult{id: id, err: err}
				return
			}
			results <- requestResult{id: id, status: response.StatusCode, body: response.Body}
		}(body, id)
	}
	close(start)
	wait.Wait()
	close(results)
	var accepted int
	for result := range results {
		var user entity.User
		query := app.Database().Table("users").Where("id = ?", result.id).Find(&user)
		persisted := query.Error == nil && user.ID == result.id && user.Email == result.id+"@example.test"
		if result.err != nil {
			if persisted {
				t.Errorf("POST user %s returned no response after persisting the row: %v", result.id, result.err)
			} else {
				t.Errorf("POST user %s failed without a response and no matching row was found: %v", result.id, result.err)
			}
			continue
		}
		if result.status != http.StatusOK {
			if persisted {
				t.Errorf("POST user %s returned %d but persisted the row; body: %s", result.id, result.status, result.body)
			} else {
				t.Errorf("POST user %s returned %d and created no row; unique valid creates must return %d, body: %s",
					result.id, result.status, http.StatusOK, result.body)
			}
			continue
		}
		var responseBody struct {
			Success bool `json:"success"`
		}
		if err := json.Unmarshal(result.body, &responseBody); err != nil {
			t.Errorf("POST user %s returned an invalid success response: %v; body: %s", result.id, err, result.body)
			continue
		}
		if !responseBody.Success {
			t.Errorf("POST user %s returned HTTP 200 with success=false; body: %s", result.id, result.body)
			continue
		}
		if !persisted {
			t.Errorf("POST user %s returned success without matching persisted state (query error: %v, row: %+v)", result.id, query.Error, user)
		}
		accepted++
	}
	if accepted != requestCount {
		t.Errorf("only %d of %d valid concurrent creates returned a consistent success response", accepted, requestCount)
	}
	for _, id := range ids {
		if query := app.Database().Table("users").Where("id = ?", id).Delete(&entity.User{}); query.Error != nil && !strings.Contains(query.Error.Error(), "record not found") {
			t.Errorf("clean up concurrent user %s: %s: %v", id, query.Message, query.Error)
		}
	}
	t.Logf("%s concurrent HTTP creates returned HTTP 200, success=true, and matching persisted rows for all %d requests", driver, accepted)
}

func verifyMigrationStartupFailure(t *testing.T, config database.Config, driver string) {
	t.Helper()

	databaseDir := t.TempDir()
	migrationDir := filepath.Join(databaseDir, "migrations", driver)
	if err := os.MkdirAll(migrationDir, 0o700); err != nil {
		t.Fatalf("create invalid %s migration directory: %v", driver, err)
	}
	extension := ".sql"
	contents := "THIS IS NOT VALID SQL;"
	if driver == database.MONGO {
		extension = ".json"
		contents = "{ invalid json"
	} else if driver == database.SCYLLA {
		extension = ".cql"
		contents = "THIS IS NOT VALID CQL;"
	}
	filename := fmt.Sprintf("999_integration_failure_%d%s", time.Now().UnixNano(), extension)
	if err := os.WriteFile(filepath.Join(migrationDir, filename), []byte(contents), 0o600); err != nil {
		t.Fatalf("write invalid %s migration: %v", driver, err)
	}
	config.DatabaseDirectory = databaseDir
	config.MigrationDirectory = "migrations"
	config.Seeder = false
	config.Migration = true
	app := newIntegrationApp(t, config)
	err := app.RunDatabaseStartupTasks()
	if err == nil {
		t.Fatalf("%s startup succeeded with an invalid migration; expected failure", driver)
	}
	if !strings.Contains(err.Error(), filename) {
		t.Fatalf("%s startup failed for an unexpected reason instead of invalid migration %q: %v", driver, filename, err)
	}
	clearRetainedStartupFailureState(t, app, config, driver, filename, "999_integration_failure"+extension)
	t.Logf("%s startup rejected an invalid migration", driver)
}

func clearRetainedStartupFailureState(t *testing.T, app *gorest.App, config database.Config, driver string, migrationNames ...string) {
	t.Helper()
	ctx, cancel := app.NewContext(config.Timeout)
	defer cancel()
	switch driver {
	case database.MONGO:
		lockCollection, _, _, err := app.Database().MongoCollection("gorest_migration_lock")
		if err != nil {
			t.Errorf("inspect failed-migration lock for %s: %v", driver, err)
			return
		}
		if _, err := lockCollection.DeleteOne(ctx, bson.M{"_id": "database_startup"}); err != nil {
			t.Errorf("clear test-owned failed-migration lock for %s: %v", driver, err)
		}
		stateCollection, _, _, err := app.Database().MongoCollection("gorest_migration_state")
		if err != nil {
			t.Errorf("inspect failed-migration state for %s: %v", driver, err)
			return
		}
		for _, name := range migrationNames {
			if _, err := stateCollection.DeleteOne(ctx, bson.M{"_id": "migration:" + name}); err != nil {
				t.Errorf("clear test-owned failed-migration state %q for %s: %v", name, driver, err)
			}
		}
	case database.SCYLLA:
		if _, _, err := app.Database().ExecCQL(ctx, "DELETE FROM gorest_migration_lock WHERE lock_name = ?", "database_startup"); err != nil {
			t.Errorf("clear test-owned failed-migration lock for %s: %v", driver, err)
		}
		for _, name := range migrationNames {
			if _, _, err := app.Database().ExecCQL(ctx, "DELETE FROM gorest_migration_state WHERE migration_name = ?", "migration:"+name); err != nil {
				t.Errorf("clear test-owned failed-migration state %q for %s: %v", name, driver, err)
			}
		}
	default:
		query := "DELETE FROM gorest_migration_lock WHERE lock_name = ?"
		if driver == database.POSTGRES {
			query = "DELETE FROM gorest_migration_lock WHERE lock_name = $1"
		} else if driver == database.SQLSERVER {
			query = "DELETE FROM gorest_migration_lock WHERE lock_name = @p1"
		} else if driver == database.ORACLE {
			query = "DELETE FROM gorest_migration_lock WHERE lock_name = :1"
		}
		if _, _, _, err := app.Database().ExecSQL(ctx, query, "database_startup"); err != nil {
			t.Errorf("clear test-owned failed-migration lock for %s: %v", driver, err)
		}
		stateQuery := "DELETE FROM gorest_migration_state WHERE migration_name = ?"
		if driver == database.POSTGRES {
			stateQuery = "DELETE FROM gorest_migration_state WHERE migration_name = $1"
		} else if driver == database.SQLSERVER {
			stateQuery = "DELETE FROM gorest_migration_state WHERE migration_name = @p1"
		} else if driver == database.ORACLE {
			stateQuery = "DELETE FROM gorest_migration_state WHERE migration_name = :1"
		}
		for _, name := range migrationNames {
			if _, _, _, err := app.Database().ExecSQL(ctx, stateQuery, "migration:"+name); err != nil {
				t.Errorf("clear test-owned failed-migration state %q for %s: %v", name, driver, err)
			}
		}
	}
}

func newUserIntegrationServer(t *testing.T, app *gorest.App) *server.Server {
	t.Helper()
	serverConfig := server.ConfigDefault
	serverConfig.Limiter = &limiter.Config{Max: 100, Expiration: time.Minute}
	testServer, _, err := server.New(&serverConfig)
	if err != nil {
		t.Fatalf("create integration HTTP server: %v", err)
	}
	userModel := model.NewUserModel(app)
	userLibrary := library.NewUserLibrary(userModel)
	userHandler := handler.NewUserHandler(app, userLibrary)
	router.RegisterUserRoutes(testServer.Group("/examples"), userHandler)
	return testServer
}

type integrationHTTPResponse struct {
	StatusCode int
	Body       []byte
}

func integrationRequest(t *testing.T, testServer *server.Server, method, path string, body []byte) integrationHTTPResponse {
	t.Helper()
	response, err := integrationRequestResult(testServer, method, path, body)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	return response
}

func integrationRequestResult(testServer *server.Server, method, path string, body []byte) (integrationHTTPResponse, error) {
	request, err := http.NewRequestWithContext(context.Background(), method, path, bytes.NewReader(body))
	if err != nil {
		return integrationHTTPResponse{}, err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := testServer.SRV.Test(request, fiber.TestConfig{Timeout: 10 * time.Second})
	if err != nil {
		return integrationHTTPResponse{}, err
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return integrationHTTPResponse{}, err
	}
	return integrationHTTPResponse{StatusCode: response.StatusCode, Body: responseBody}, nil
}

func assertDriverUserAbsent(t *testing.T, app *gorest.App, id string) {
	t.Helper()
	var user entity.User
	if query := app.Database().Table("users").Where("id = ?", id).Find(&user); query.Error == nil {
		t.Errorf("user %s unexpectedly exists after rejected request: %+v", id, user)
	}
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
