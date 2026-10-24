package e2e_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/elsyahtech/gorest/database"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func e2eDatabaseConfig(t *testing.T, dataDir, name string) database.Config {
	t.Helper()

	driver := os.Getenv("GOREST_E2E_DRIVER")
	if driver == "" {
		driver = database.SQLITE
	}
	config := database.Config{
		Driver:                 driver,
		Name:                   name,
		Host:                   os.Getenv("GOREST_DRIVER_HOST"),
		Port:                   os.Getenv("GOREST_DRIVER_PORT"),
		Username:               os.Getenv("GOREST_DRIVER_USERNAME"),
		Password:               os.Getenv("GOREST_DRIVER_PASSWORD"),
		LibDir:                 os.Getenv("GOREST_DRIVER_LIB_DIR"),
		DatabaseDirectory:      dataDir,
		MigrationDirectory:     "migrations",
		SeederDirectory:        "seeders",
		Timeout:                30 * time.Second,
		ConnectionMaxLifetime:  time.Minute,
		MaxOpenConnections:     5,
		MaxIdleConnections:     2,
		Migration:              true,
		Disable:                false,
		TrustServerCertificate: true,
	}
	if driver == database.SQLITE {
		config.Name = name
		config.MaxOpenConnections = 1
		config.MaxIdleConnections = 1
		config.DatabaseDirectory = dataDir
		return config
	}
	config.Name = os.Getenv("GOREST_DRIVER_DATABASE")
	if config.Name == "" {
		config.Name = name
	}
	if driver != database.MONGO && driver != database.SCYLLA {
		config.TLS = "disable"
	}
	return config
}

func e2eMigrationDirectory() string {
	driver := os.Getenv("GOREST_E2E_DRIVER")
	if driver == "" {
		driver = database.SQLITE
	}
	return driver
}

func e2eSQLPlaceholder(index int) string {
	switch os.Getenv("GOREST_E2E_DRIVER") {
	case database.POSTGRES:
		return fmt.Sprintf("$%d", index)
	case database.SQLSERVER:
		return fmt.Sprintf("@p%d", index)
	case database.ORACLE:
		return fmt.Sprintf(":%d", index)
	default:
		return "?"
	}
}

func resetSQLTestTables(t *testing.T, app *gorest.App) {
	t.Helper()
	ctx, cancel := app.NewContext(30 * time.Second)
	defer cancel()

	switch os.Getenv("GOREST_E2E_DRIVER") {
	case database.MONGO:
		_, _, err := app.Database().RawQueryMongo(ctx, func(ctx context.Context, client *mongo.Client) error {
			for _, collectionName := range []string{"users", "departments", "orders"} {
				if _, err := client.Database(os.Getenv("GOREST_DRIVER_DATABASE")).
					Collection(collectionName).DeleteMany(ctx, bson.D{}); err != nil {
					return fmt.Errorf("clear collection %q: %w", collectionName, err)
				}
			}
			return nil
		})
		if err != nil {
			t.Fatalf("clear MongoDB E2E collections: %v", err)
		}
		return
	case database.SCYLLA:
		for _, table := range []string{"users", "departments", "orders"} {
			if _, _, err := app.Database().ExecCQL(ctx, "TRUNCATE "+table); err != nil {
				t.Fatalf("clear ScyllaDB E2E table %q: %v", table, err)
			}
		}
		return
	}

	for _, table := range []string{"orders", "users", "departments"} {
		if _, _, _, err := app.Database().ExecSQL(ctx, "DELETE FROM "+table); err != nil {
			t.Fatalf("clear shared E2E table %q: %v", table, err)
		}
	}
}

func ensureE2EUserTimestamp(user *entity.User) {
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now().UTC()
	}
}

func ensureE2EDepartmentTimestamp(department *entity.Department) {
	if department.CreatedAt.IsZero() {
		department.CreatedAt = time.Now().UTC()
	}
}

func copyE2EMigration(t *testing.T, dataDir, sourceName string) {
	t.Helper()
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate E2E database config source")
	}
	repositoryRoot := filepath.Clean(filepath.Join(filepath.Dir(sourceFile), "../.."))
	driver := e2eMigrationDirectory()
	actualSourceName := sourceName
	switch driver {
	case database.MONGO:
		actualSourceName = map[string]string{
			"001_create_departments_table.sql": "001_create_departments_collection.json",
			"002_create_users_table.sql":       "002_create_users_collection.json",
			"003_create_orders_table.sql":      "003_create_orders_collection.json",
		}[sourceName]
	case database.SCYLLA:
		actualSourceName = strings.TrimSuffix(sourceName, ".sql") + ".cql"
	}
	if actualSourceName == "" {
		t.Fatalf("no %s E2E migration matches %q", driver, sourceName)
	}
	source := filepath.Join(repositoryRoot, "app", "database", "migrations", driver, actualSourceName)
	destination := filepath.Join(dataDir, "migrations", e2eMigrationDirectory(), filepath.Base(source))
	migration, err := os.ReadFile(source)
	if err != nil {
		t.Fatalf("read %s migration: %v", e2eMigrationDirectory(), err)
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		t.Fatalf("create E2E migration directory: %v", err)
	}
	if err := os.WriteFile(destination, migration, 0o600); err != nil {
		t.Fatalf("write E2E migration: %v", err)
	}
}
