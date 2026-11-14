//go:build integration

package integration_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/elsyahtech/gorest/database"
	"github.com/joho/godotenv"
)

type sqlServerTestConfig struct {
	host              string
	port              string
	database          string
	username          string
	password          string
	tls               string
	caPath            string
	databaseDirectory string
}

func loadSQLServerTestConfig() (sqlServerTestConfig, error) {
	values, err := godotenv.Read("../.env")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return sqlServerTestConfig{}, fmt.Errorf("read .env: %w", err)
	}
	value := func(key string) string {
		if current, ok := os.LookupEnv(key); ok {
			return current
		}
		return values[key]
	}

	if strings.EqualFold(strings.TrimSpace(value("GOREST_APP_ENV")), "production") {
		return sqlServerTestConfig{}, errors.New("SQL Server integration tests refuse to run with GOREST_APP_ENV=production")
	}
	if strings.ToLower(strings.TrimSpace(value("GOREST_DB_DRIVER"))) != database.SQLSERVER {
		return sqlServerTestConfig{}, errors.New("GOREST_DB_DRIVER=sqlserver is required for SQL Server integration tests")
	}
	if strings.EqualFold(strings.TrimSpace(value("GOREST_DB_DISABLE")), "true") {
		return sqlServerTestConfig{}, errors.New("GOREST_DB_DISABLE=false is required for SQL Server integration tests")
	}
	if strings.EqualFold(strings.TrimSpace(value("GOREST_DB_TRUST_SERVER_CERTIFICATE")), "true") {
		return sqlServerTestConfig{}, errors.New("SQL Server integration tests require GOREST_DB_TRUST_SERVER_CERTIFICATE=false to verify the server certificate")
	}

	config := sqlServerTestConfig{
		host:     strings.TrimSpace(value("GOREST_DB_HOST")),
		port:     strings.TrimSpace(value("GOREST_DB_PORT")),
		database: strings.TrimSpace(value("GOREST_DB_NAME")),
		username: strings.TrimSpace(value("GOREST_DB_USER")),
		password: strings.TrimSpace(value("GOREST_DB_PASSWORD")),
		tls:      strings.ToLower(strings.TrimSpace(value("GOREST_DB_TLS"))),
		caPath:   strings.TrimSpace(value("GOREST_DB_CA_PATH")),
	}
	if config.port == "" {
		config.port = database.SQLSERVER_PORT
	}
	if config.tls != database.TLSMandatory && config.tls != database.TLSStrict {
		return sqlServerTestConfig{}, errors.New("GOREST_DB_TLS must be mandatory or strict for SQL Server integration tests")
	}
	for key, configured := range map[string]string{
		"GOREST_DB_HOST":     config.host,
		"GOREST_DB_NAME":     config.database,
		"GOREST_DB_USER":     config.username,
		"GOREST_DB_PASSWORD": config.password,
	} {
		if configured == "" {
			return sqlServerTestConfig{}, fmt.Errorf("%s is required for SQL Server integration tests", key)
		}
	}
	switch strings.ToLower(config.database) {
	case "master", "model", "msdb", "tempdb":
		return sqlServerTestConfig{}, errors.New("SQL Server integration tests refuse to run against a system database")
	}
	if strings.TrimSpace(value("GOREST_TEST_SQLSERVER_CONFIRM_DATABASE")) != config.database {
		return sqlServerTestConfig{}, errors.New("set GOREST_TEST_SQLSERVER_CONFIRM_DATABASE to the exact SQL Server test database name to allow migrations, seeders, and test data writes")
	}

	root, err := filepath.Abs("..")
	if err != nil {
		return sqlServerTestConfig{}, fmt.Errorf("resolve repository root for SQL Server paths: %w", err)
	}
	if config.caPath != "" && !filepath.IsAbs(config.caPath) {
		config.caPath = filepath.Join(root, config.caPath)
	}
	if config.caPath != "" {
		if _, err := os.Stat(config.caPath); err != nil {
			return sqlServerTestConfig{}, fmt.Errorf("access configured SQL Server CA bundle: %w", err)
		}
	}
	config.databaseDirectory = filepath.Join(root, "app", "database")
	return config, nil
}

func (config sqlServerTestConfig) appEnvironment(serverPort, logDirectory, signatureKey string) map[string]string {
	return map[string]string{
		"GOREST_APP_ENV":                      "testing",
		"GOREST_APP_DEBUG":                    "false",
		"GOREST_APP_CHECK_FOR_UPDATES":        "false",
		"GOREST_SERVER_HOST":                  "127.0.0.1",
		"GOREST_SERVER_PORT":                  serverPort,
		"GOREST_SERVER_LISTENER_NETWORK":      "tcp4",
		"GOREST_SERVER_CERT_FILE":             "",
		"GOREST_SERVER_CERT_KEY_FILE":         "",
		"GOREST_SERVER_CERT_CLIENT_FILE":      "",
		"GOREST_SERVER_ENABLE_PREFORK":        "false",
		"GOREST_SERVER_CORS_ORIGINS":          "http://127.0.0.1",
		"GOREST_DB_DISABLE":                   "false",
		"GOREST_DB_DRIVER":                    database.SQLSERVER,
		"GOREST_DB_HOST":                      config.host,
		"GOREST_DB_PORT":                      config.port,
		"GOREST_DB_NAME":                      config.database,
		"GOREST_DB_USER":                      config.username,
		"GOREST_DB_PASSWORD":                  config.password,
		"GOREST_DB_TLS":                       config.tls,
		"GOREST_DB_CA_PATH":                   config.caPath,
		"GOREST_DB_TRUST_SERVER_CERTIFICATE":  "false",
		"GOREST_DB_DIRECTORY":                 config.databaseDirectory,
		"GOREST_DB_MAX_IDLE_CONNECTIONS":      "0",
		"GOREST_DB_MIGRATION":                 "true",
		"GOREST_DB_SEEDER":                    "true",
		"GOREST_DB_ALLOW_INSECURE_PRODUCTION": "false",
		"GOREST_REDIS_DISABLE":                "true",
		"GOREST_FILTER_SIGNATURE_KEY":         signatureKey,
		"GOREST_FILTER_COOKIE_SECURE":         "false",
		"GOREST_MONITORING_ALLOWED_IPS":       "127.0.0.1,::1",
		"GOREST_MONITORING_DISABLE":           "false",
		"GOREST_LOG_DIR":                      logDirectory,
		"GOREST_LOG_FILE":                     "integration-test.log",
		"GOREST_LOG_CALLER":                   "false",
		"GOREST_LOG_MAX_FILE_SIZE":            "1",
		"GOREST_LOG_MAX_BACKUPS":              "1",
		"GOREST_LOG_MAX_AGE":                  "1",
	}
}
