package drivers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/elsyahtech/gorest"
	appconfig "github.com/elsyahtech/gorest-api-starter/app/config"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/elsyahtech/gorest-api-starter/app/handler"
	"github.com/elsyahtech/gorest-api-starter/app/library"
	"github.com/elsyahtech/gorest-api-starter/app/model"
	"github.com/elsyahtech/gorest-api-starter/app/router"
	"github.com/elsyahtech/gorest-api-starter/test/mysqltls"
	"github.com/elsyahtech/gorest/database"
	"github.com/elsyahtech/gorest/filter"
	gorestlog "github.com/elsyahtech/gorest/log"
	"github.com/elsyahtech/gorest/monitoring"
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
	if driver == database.SQLITE {
		verifySQLiteForeignKeyEnforcement(t, app)
	}
	if driver == database.SQLITE || driver == database.POSTGRES || driver == database.MYSQL || driver == database.ORACLE || driver == database.SQLSERVER {
		verifyNativeSQLRoutes(t, app, driver)
	}
	if driver != database.MONGO && driver != database.SCYLLA {
		verifyTransactionRollback(t, app, driver)
		verifyMidStepRollback(t, app, driver)
	} else {
		t.Logf("%s transaction atomicity: not applicable to the starter's GOREST ORM transaction API", driver)
	}
	verifyConflictAndInvalidRequests(t, app, driver)
	verifyConcurrentRequests(t, app, driver)
	verifyConnectionStartupFailure(t, driver)
	verifyMigrationStartupFailure(t, config, driver)
}

func TestMySQLStorageEndpointIntegration(t *testing.T) {
	driver := strings.ToLower(strings.TrimSpace(os.Getenv("GOREST_DRIVER_INTEGRATION")))
	if driver != database.MYSQL {
		t.Skip("set GOREST_DRIVER_INTEGRATION=mysql to run the MySQL storage endpoint integration test")
	}

	config, _ := integrationDatabaseConfig(t, driver)
	verifyMySQLStorageEndpoint(t, config)
}

func verifyMySQLStorageEndpoint(t *testing.T, config database.Config) {
	t.Helper()

	config = database.Run(config)
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	storageDatabase, message, err := database.New(ctx, &config, "UTC")
	if err != nil {
		t.Fatalf("open MySQL storage endpoint database connection (%s): %v", message, err)
	}
	t.Cleanup(func() {
		closeCtx, closeCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer closeCancel()
		if message, _, err := storageDatabase.Close(closeCtx, &config); err != nil {
			t.Errorf("close MySQL storage endpoint database connection (%s): %v", message, err)
		}
	})

	httpServer, _, err := server.New(nil)
	if err != nil {
		t.Fatalf("create HTTP server for MySQL storage endpoint: %v", err)
	}
	monitoring.Register(httpServer.Group("/gorest"), config, storageDatabase)

	beforeBytes := requestMySQLStorageMetric(t, httpServer)

	tableName := fmt.Sprintf("gorest_storage_metric_%d", time.Now().UnixNano())
	if _, err := storageDatabase.SQL.ExecContext(ctx, fmt.Sprintf(
		"CREATE TABLE `%s` (id INT NOT NULL PRIMARY KEY, payload VARCHAR(4096) NOT NULL) ENGINE=InnoDB",
		tableName,
	)); err != nil {
		t.Fatalf("create temporary MySQL storage metric table: %v", err)
	}
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()
		if _, err := storageDatabase.SQL.ExecContext(cleanupCtx, fmt.Sprintf("DROP TABLE IF EXISTS `%s`", tableName)); err != nil {
			t.Errorf("drop temporary MySQL storage metric table: %v", err)
		}
	})

	tx, err := storageDatabase.SQL.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin MySQL storage metric inserts: %v", err)
	}
	payload := strings.Repeat("x", 4000)
	for firstID := 1; firstID <= 500; firstID += 50 {
		const batchSize = 50
		values := strings.TrimSuffix(strings.Repeat("(?, ?),", batchSize), ",")
		args := make([]any, 0, batchSize*2)
		for id := firstID; id < firstID+batchSize; id++ {
			args = append(args, id, payload)
		}
		if _, err := tx.ExecContext(ctx,
			fmt.Sprintf("INSERT INTO `%s` (id, payload) VALUES %s", tableName, values),
			args...,
		); err != nil {
			_ = tx.Rollback()
			t.Fatalf("insert MySQL storage metric rows %d-%d of 500: %v", firstID, firstID+batchSize-1, err)
		}
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("commit MySQL storage metric inserts: %v", err)
	}

	afterBytes := requestMySQLStorageMetric(t, httpServer)
	t.Logf("mysql storage estimate before 500 rows: %d bytes; after: %d bytes; delta: %d bytes",
		beforeBytes, afterBytes, afterBytes-beforeBytes)
}

func requestMySQLStorageMetric(t *testing.T, httpServer *server.Server) int64 {
	t.Helper()

	response, err := httpServer.SRV.Test(httptest.NewRequest(http.MethodGet, "/gorest/storage", nil))
	if err != nil {
		t.Fatalf("request MySQL storage endpoint: %v", err)
	}
	body, readErr := io.ReadAll(response.Body)
	closeErr := response.Body.Close()
	if readErr != nil {
		t.Fatalf("read MySQL storage endpoint response: %v", readErr)
	}
	if closeErr != nil {
		t.Fatalf("close MySQL storage endpoint response: %v", closeErr)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("MySQL storage endpoint status = %d, want %d; response: %s",
			response.StatusCode, http.StatusOK, body)
	}

	var result struct {
		Status   string `json:"status"`
		Database struct {
			Driver      string `json:"driver"`
			Scope       string `json:"scope"`
			Measurement string `json:"measurement"`
			UsedBytes   int64  `json:"used_bytes"`
		} `json:"database"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("decode MySQL storage endpoint response: %v; response: %s", err, body)
	}
	if result.Status != "available" {
		t.Fatalf("MySQL storage endpoint status field = %q, want available; response: %s", result.Status, body)
	}
	if result.Database.Driver != database.MYSQL ||
		result.Database.Scope != "configured_database" ||
		result.Database.Measurement != "table_data_and_index_size" ||
		result.Database.UsedBytes < 0 {
		t.Fatalf("unexpected MySQL storage metric: %+v", result.Database)
	}
	t.Logf("mysql storage endpoint passed; database estimate=%d bytes (%s)",
		result.Database.UsedBytes, result.Database.Measurement)
	return result.Database.UsedBytes
}

func verifyNativeSQLRoutes(t *testing.T, app *gorest.App, driver string) {
	t.Helper()
	testServer := newUserIntegrationServer(t, app)
	id := "driver-native-sql-" + driver
	user := entity.User{
		ID: id, Email: id + "@example.test", Password: "native-sql-test-password",
		DepartmentID: "dept-01", IsActive: true,
	}
	body, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("encode %s native SQL user: %v", driver, err)
	}
	path := "/examples/native-queries/users/sql"
	created := integrationRequest(t, testServer, http.MethodPost, path, body)
	if created.StatusCode < http.StatusOK || created.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("%s native SQL create status = %d, want 2xx; body: %s", driver, created.StatusCode, created.Body)
	}

	updated := user
	updated.Email = "updated-" + id + "@example.test"
	updated.Password = "updated-native-sql-test-password"
	updatedBody, err := json.Marshal(updated)
	if err != nil {
		t.Fatalf("encode %s native SQL update: %v", driver, err)
	}
	response := integrationRequest(t, testServer, http.MethodPut, path+"/"+id, updatedBody)
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("%s native SQL update status = %d, want 2xx; body: %s", driver, response.StatusCode, response.Body)
	}

	list := integrationRequest(t, testServer, http.MethodGet, path, nil)
	if list.StatusCode != http.StatusOK || !bytes.Contains(list.Body, []byte(updated.Email)) {
		t.Fatalf("%s native SQL read status/body do not include updated user: status=%d body=%s", driver, list.StatusCode, list.Body)
	}

	deleted := integrationRequest(t, testServer, http.MethodDelete, path+"/"+id, nil)
	if deleted.StatusCode < http.StatusOK || deleted.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("%s native SQL delete status = %d, want 2xx; body: %s", driver, deleted.StatusCode, deleted.Body)
	}
	var persisted entity.User
	query := app.Database().Table("users").Where("id = ?", id).Find(&persisted)
	if query.Error == nil {
		t.Fatalf("%s native SQL delete left user persisted: %+v", driver, persisted)
	}
	t.Logf("%s native SQL route create/read/update/delete passed", driver)
}

func TestDatabaseDriverStartupFailsWithoutDatabase(t *testing.T) {
	driver := strings.ToLower(strings.TrimSpace(os.Getenv("GOREST_DRIVER_LIFECYCLE_DRIVER")))
	if driver == "" {
		t.Skip("set GOREST_DRIVER_LIFECYCLE_DRIVER to verify connection failure during startup")
	}
	verifyConnectionStartupFailure(t, driver)
}

func TestDatabaseDriverProcessFailsOnConnection(t *testing.T) {
	driver := strings.ToLower(strings.TrimSpace(os.Getenv("GOREST_DRIVER_CONNECTION_FAILURE")))
	if driver == "" {
		t.Skip("internal subprocess check")
	}
	startWithUnavailableDatabase(t, driver)
}

func verifyConnectionStartupFailure(t *testing.T, driver string) {
	t.Helper()

	config, _ := integrationDatabaseConfig(t, driver)
	config.Migration = false
	config.Seeder = false
	if driver == database.SQLITE {
		t.Log("sqlite startup connection failure: not applicable; SQLite is embedded. SQLite migration startup failure is covered by the migration failure scenario.")
		return
	}
	childValues := map[string]string{"GOREST_DRIVER_CONNECTION_FAILURE": driver}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve unavailable %s database port: %v", driver, err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		t.Fatalf("release unavailable %s database port: %v", driver, err)
	}
	childValues["GOREST_DRIVER_FAILURE_PORT"] = fmt.Sprint(port)

	executable, err := os.Executable()
	if err != nil {
		t.Fatalf("locate test executable for %s connection failure check: %v", driver, err)
	}
	command := exec.Command(executable, "-test.run=^TestDatabaseDriverProcessFailsOnConnection$", "-test.v")
	command.Env = integrationChildEnvironment(childValues)
	output, err := command.CombinedOutput()
	if err == nil {
		t.Fatalf("%s application startup succeeded with an unavailable database; output: %s", driver, output)
	}
	if !strings.Contains(string(output), "failed to run database") {
		t.Fatalf("%s startup exited for an unexpected reason; output: %s", driver, output)
	}
	if strings.Contains(string(output), "database initialized successfully") {
		t.Fatalf("%s startup reported database initialization success despite connection failure; output: %s", driver, output)
	}
	t.Logf("%s startup rejected an unavailable database", driver)
}

func startWithUnavailableDatabase(t *testing.T, driver string) {
	t.Helper()

	config, _ := integrationDatabaseConfig(t, driver)
	config.Migration = false
	config.Seeder = false
	config.Port = os.Getenv("GOREST_DRIVER_FAILURE_PORT")
	logConfig := gorestlog.ConfigDefault
	logConfig.LogDir = filepath.Join(t.TempDir(), "logs")
	app := gorest.New(gorest.Config{AppName: "gorest-api-starter-driver-startup-failure", Environment: gorest.Testing, Timezone: "UTC"})
	app.Use(gorestlog.Run(logConfig), database.Run(config))
	t.Fatal("application startup unexpectedly succeeded with an unavailable database")
}

func integrationChildEnvironment(overrides map[string]string) []string {
	environment := make([]string, 0, len(os.Environ())+len(overrides))
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if _, overridden := overrides[key]; !overridden {
			environment = append(environment, entry)
		}
	}
	for key, value := range overrides {
		environment = append(environment, key+"="+value)
	}
	return environment
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
		Disable:                false,
		Migration:              true,
		Seeder:                 true,
		TrustServerCertificate: true,
	}
	databaseName := envOr("GOREST_DRIVER_DATABASE", "gorest_driver_it")

	switch driver {
	case database.SQLITE:
		// SQLite has a single-writer model; serialize DB work in the test pool
		// while keeping all HTTP requests concurrent and foreign keys enabled.
		config.MaxOpenConnections = 1
		config.MaxIdleConnections = 1
		config.ConnectionMaxLifetime = 0
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
		config.TLS = envOr("GOREST_DRIVER_TLS", "disable")
		config.SSLCAPath = os.Getenv("GOREST_DRIVER_SSL_CA_PATH")
		config.SSLCertPath = os.Getenv("GOREST_DRIVER_SSL_CERT_PATH")
		config.SSLKeyPath = os.Getenv("GOREST_DRIVER_SSL_KEY_PATH")
		if config.TLS != "disable" && config.TLS != "true" {
			if err := mysqltls.RegisterProfile(
				config.TLS,
				envOr("GOREST_DRIVER_TLS_SERVER_NAME", config.Host),
				config.SSLCAPath,
				config.SSLCertPath,
				config.SSLKeyPath,
			); err != nil {
				t.Fatalf("configure MySQL TLS: %v", err)
			}
		}
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
		Environment: gorest.Testing,
		Timezone:    "UTC",
	})
	app.Use(
		gorestlog.Run(logConfig),
		database.Run(config),
		filter.Run(driverIntegrationFilterConfig()),
	)
	if config.Driver == database.SQLITE {
		if err := appconfig.EnableSQLiteForeignKeys(app); err != nil {
			t.Fatalf("enable SQLite foreign keys: %v", err)
		}
	}
	t.Cleanup(func() {
		if message, _, err := app.Close(); err != nil {
			t.Errorf("close %s integration app (%s): %v", config.Driver, message, err)
		} else {
			t.Logf("%s application shutdown passed", config.Driver)
		}
	})
	return app
}

func driverIntegrationFilterConfig() filter.Config {
	signatureKey := "gorest-api-starter-driver-integration-test-key"
	config := filter.ConfigDefault
	config.SignatureKey = &signatureKey
	config.CookieSecure = false
	config.LoggerConfig.DisableLogging = true
	return config
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

func verifySQLiteForeignKeyEnforcement(t *testing.T, app *gorest.App) {
	t.Helper()

	const invalidOrderID = "sqlite-foreign-key-invalid-user"
	invalidOrder := entity.Order{
		ID:          invalidOrderID,
		UserID:      "sqlite-user-does-not-exist",
		OrderNumber: invalidOrderID,
		Status:      "PENDING",
		TotalAmount: 1,
		IsActive:    true,
		CreatedAt:   time.Now().UTC(),
	}
	query := app.Database().Table("orders").Create(&invalidOrder)
	if query.Error == nil {
		t.Fatal("SQLite accepted an order insert referencing a missing user; want foreign-key constraint failure")
	}

	var persisted entity.Order
	query = app.Database().Table("orders").Where("id = ?", invalidOrderID).Find(&persisted)
	if query.Error == nil {
		t.Fatalf("SQLite persisted an order insert with a missing user: %+v", persisted)
	}

	var user entity.User
	query = app.Database().Table("users").Find(&user)
	if query.Error != nil {
		t.Fatalf("read a seeded user for the SQLite foreign-key update check: %s: %v",
			query.Message, query.Error)
	}
	validOrder := entity.Order{
		ID:          "sqlite-foreign-key-update",
		UserID:      user.ID,
		OrderNumber: "sqlite-foreign-key-update",
		Status:      "PENDING",
		TotalAmount: 1,
		IsActive:    true,
		CreatedAt:   time.Now().UTC(),
	}
	query = app.Database().Table("orders").Create(&validOrder)
	if query.Error != nil {
		t.Fatalf("create SQLite order for foreign-key update check: %s: %v", query.Message, query.Error)
	}
	validUserID := validOrder.UserID
	validOrder.UserID = "sqlite-user-does-not-exist"
	query = app.Database().Table("orders").
		Where("id = ?", validOrder.ID).Update(&validOrder)
	if query.Error == nil {
		t.Fatal("SQLite accepted an order update referencing a missing user; want foreign-key constraint failure")
	}
	query = app.Database().Table("orders").Where("id = ?", validOrder.ID).Find(&persisted)
	if query.Error != nil || persisted.UserID != validUserID {
		t.Fatalf("SQLite order after rejected foreign-key update = %+v, query error = %v; want user %q",
			persisted, query.Error, validUserID)
	}
	t.Log("SQLite foreign-key constraints rejected invalid order inserts and updates")
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
	userLibrary := library.NewUserLibrary(app, userModel)
	userHandler := handler.NewUserHandler(app, userLibrary)
	router.RegisterUserRoutes(app, testServer.Group("/examples"), userHandler)
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
	response, err := testServer.SRV.Test(request, fiber.TestConfig{Timeout: 30 * time.Second})
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
