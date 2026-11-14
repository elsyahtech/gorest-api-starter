//go:build integration

package integration_test

import (
	"path/filepath"
	"testing"

	"github.com/elsyahtech/gorest/database"
)

type sqliteTestConfig struct {
	database            string
	databaseDirectory   string
	migrationsDirectory string
	seedersDirectory    string
}

func newSQLiteTestConfig(t *testing.T) sqliteTestConfig {
	t.Helper()

	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("resolve repository root for SQLite integration assets: %v", err)
	}
	databaseDirectory := t.TempDir()
	migrationsDirectory, err := filepath.Rel(
		databaseDirectory,
		filepath.Join(root, "app", "database", "migrations"),
	)
	if err != nil {
		t.Fatalf("resolve SQLite migrations directory: %v", err)
	}
	seedersDirectory, err := filepath.Rel(
		databaseDirectory,
		filepath.Join(root, "app", "database", "seeders"),
	)
	if err != nil {
		t.Fatalf("resolve SQLite seeders directory: %v", err)
	}

	return sqliteTestConfig{
		database:            "integration_test",
		databaseDirectory:   databaseDirectory,
		migrationsDirectory: migrationsDirectory,
		seedersDirectory:    seedersDirectory,
	}
}

func (config sqliteTestConfig) databaseEnvironment() map[string]string {
	return map[string]string{
		"GOREST_DB_DISABLE":              "false",
		"GOREST_DB_DRIVER":               database.SQLITE,
		"GOREST_DB_NAME":                 config.database,
		"GOREST_DB_DIRECTORY":            config.databaseDirectory,
		"GOREST_DB_MIGRATION":            "true",
		"GOREST_DB_SEEDER":               "true",
		"GOREST_DB_MIGRATIONS_DIRECTORY": config.migrationsDirectory,
		"GOREST_DB_SEEDERS_DIRECTORY":    config.seedersDirectory,
	}
}

func (config sqliteTestConfig) appEnvironment(serverPort, logDirectory, signatureKey string) map[string]string {
	environment := config.databaseEnvironment()
	for key, value := range map[string]string{
		"GOREST_APP_ENV":                 "testing",
		"GOREST_APP_DEBUG":               "false",
		"GOREST_APP_CHECK_FOR_UPDATES":   "false",
		"GOREST_SERVER_HOST":             "127.0.0.1",
		"GOREST_SERVER_PORT":             serverPort,
		"GOREST_SERVER_LISTENER_NETWORK": "tcp4",
		"GOREST_SERVER_CERT_FILE":        "",
		"GOREST_SERVER_CERT_KEY_FILE":    "",
		"GOREST_SERVER_CERT_CLIENT_FILE": "",
		"GOREST_SERVER_ENABLE_PREFORK":   "false",
		"GOREST_SERVER_CORS_ORIGINS":     "http://127.0.0.1",
		"GOREST_REDIS_DISABLE":           "true",
		"GOREST_FILTER_SIGNATURE_KEY":    signatureKey,
		"GOREST_FILTER_COOKIE_SECURE":    "false",
		"GOREST_MONITORING_ALLOWED_IPS":  "127.0.0.1,::1",
		"GOREST_MONITORING_DISABLE":      "false",
		"GOREST_LOG_DIR":                 logDirectory,
		"GOREST_LOG_FILE":                "integration-test.log",
		"GOREST_LOG_CALLER":              "false",
		"GOREST_LOG_MAX_FILE_SIZE":       "1",
		"GOREST_LOG_MAX_BACKUPS":         "1",
		"GOREST_LOG_MAX_AGE":             "1",
	} {
		environment[key] = value
	}
	return environment
}

func applySQLiteTestConfig(t *testing.T, config sqliteTestConfig) {
	t.Helper()
	for key, value := range config.databaseEnvironment() {
		t.Setenv(key, value)
	}
}
