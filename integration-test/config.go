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

type postgresTestConfig struct {
	host              string
	port              string
	database          string
	username          string
	password          string
	tls               string
	caPath            string
	certPath          string
	keyPath           string
	databaseDirectory string
}

func loadPostgreSQLTestConfig() (postgresTestConfig, error) {
	values, err := godotenv.Read("../.env")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return postgresTestConfig{}, fmt.Errorf("read .env: %w", err)
	}
	value := func(key string) string {
		if current, ok := os.LookupEnv(key); ok {
			return current
		}
		return values[key]
	}

	if strings.EqualFold(strings.TrimSpace(value("GOREST_APP_ENV")), "production") {
		return postgresTestConfig{}, errors.New("PostgreSQL integration tests refuse to run with GOREST_APP_ENV=production")
	}
	if strings.ToLower(strings.TrimSpace(value("GOREST_DB_DRIVER"))) != database.POSTGRES {
		return postgresTestConfig{}, errors.New("GOREST_DB_DRIVER=postgres is required for PostgreSQL integration tests")
	}
	if strings.EqualFold(strings.TrimSpace(value("GOREST_DB_DISABLE")), "true") {
		return postgresTestConfig{}, errors.New("GOREST_DB_DISABLE=false is required for PostgreSQL integration tests")
	}
	certPath := strings.TrimSpace(value("GOREST_DB_CERT_PATH"))
	keyPath := strings.TrimSpace(value("GOREST_DB_KEY_PATH"))
	if certPath == "" || keyPath == "" {
		return postgresTestConfig{}, errors.New("PostgreSQL mTLS integration tests require both GOREST_DB_CERT_PATH and GOREST_DB_KEY_PATH")
	}

	config := postgresTestConfig{
		host:     strings.TrimSpace(value("GOREST_DB_HOST")),
		port:     strings.TrimSpace(value("GOREST_DB_PORT")),
		database: strings.TrimSpace(value("GOREST_DB_NAME")),
		username: value("GOREST_DB_USER"),
		password: value("GOREST_DB_PASSWORD"),
		tls:      strings.ToLower(strings.TrimSpace(value("GOREST_DB_TLS"))),
		caPath:   strings.TrimSpace(value("GOREST_DB_CA_PATH")),
		certPath: certPath,
		keyPath:  keyPath,
	}
	if config.port == "" {
		config.port = database.POSTGRES_PORT
	}
	switch config.tls {
	case "require", "verify-ca", "verify-full":
	default:
		return postgresTestConfig{}, fmt.Errorf("GOREST_DB_TLS must be require, verify-ca, or verify-full; configured mode is %q", config.tls)
	}
	for key, configured := range map[string]string{
		"GOREST_DB_HOST":     config.host,
		"GOREST_DB_NAME":     config.database,
		"GOREST_DB_USER":     config.username,
		"GOREST_DB_PASSWORD": config.password,
	} {
		if configured == "" {
			return postgresTestConfig{}, fmt.Errorf("%s is required for PostgreSQL integration tests", key)
		}
	}

	confirmedDatabase := strings.TrimSpace(value("GOREST_TEST_POSTGRES_CONFIRM_DATABASE"))
	allowedDatabase := strings.TrimSpace(value("GOREST_TEST_POSTGRES_ALLOW_NON_TEST_DATABASE"))
	if !strings.HasSuffix(strings.ToLower(config.database), "_test") &&
		!strings.HasSuffix(strings.ToLower(config.database), "_testing") &&
		(confirmedDatabase != config.database || allowedDatabase != config.database) {
		return postgresTestConfig{}, errors.New("PostgreSQL integration tests require a _test/_testing database or matching GOREST_TEST_POSTGRES_CONFIRM_DATABASE and GOREST_TEST_POSTGRES_ALLOW_NON_TEST_DATABASE")
	}
	if confirmedDatabase != config.database {
		return postgresTestConfig{}, errors.New("set GOREST_TEST_POSTGRES_CONFIRM_DATABASE to the exact PostgreSQL database name to allow migrations, seeders, and test data writes")
	}

	root, err := filepath.Abs("..")
	if err != nil {
		return postgresTestConfig{}, fmt.Errorf("resolve repository root for PostgreSQL paths: %w", err)
	}
	for _, path := range []*string{&config.caPath, &config.certPath, &config.keyPath} {
		if *path != "" && !filepath.IsAbs(*path) {
			*path = filepath.Join(root, *path)
		}
	}
	config.databaseDirectory = filepath.Join(root, "app", "database")
	return config, nil
}

func (config postgresTestConfig) appEnvironment(serverPort, logDirectory, signatureKey string) map[string]string {
	return map[string]string{
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
		"GOREST_DB_DISABLE":              "false",
		"GOREST_DB_DRIVER":               database.POSTGRES,
		"GOREST_DB_HOST":                 config.host,
		"GOREST_DB_PORT":                 config.port,
		"GOREST_DB_NAME":                 config.database,
		"GOREST_DB_USER":                 config.username,
		"GOREST_DB_PASSWORD":             config.password,
		"GOREST_DB_TLS":                  config.tls,
		"GOREST_DB_CA_PATH":              config.caPath,
		"GOREST_DB_CERT_PATH":            config.certPath,
		"GOREST_DB_KEY_PATH":             config.keyPath,
		"GOREST_DB_DIRECTORY":            config.databaseDirectory,
		"GOREST_DB_MAX_IDLE_CONNECTIONS": "0",
		"GOREST_DB_MIGRATION":            "true",
		"GOREST_DB_SEEDER":               "true",
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
	}
}
