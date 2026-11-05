//go:build integration

package integration_test

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type mysqlTestConfig struct {
	host     string
	port     string
	database string
	username string
	password string
	tls      string
	caPath   string
	certPath string
	keyPath  string
}

func loadMySQLTestConfig() (mysqlTestConfig, error) {
	values, err := godotenv.Read("../.env")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return mysqlTestConfig{}, fmt.Errorf("read .env: %w", err)
	}

	value := func(testKey, appKey string) string {
		if current := os.Getenv(testKey); current != "" {
			return current
		}
		if configured := values[testKey]; configured != "" {
			return configured
		}
		if current := os.Getenv(appKey); current != "" {
			return current
		}
		return values[appKey]
	}

	usesApplicationConfig := false
	for _, key := range []string{
		"GOREST_TEST_MYSQL_HOST",
		"GOREST_TEST_MYSQL_DATABASE",
		"GOREST_TEST_MYSQL_USER",
		"GOREST_TEST_MYSQL_PASSWORD",
	} {
		if os.Getenv(key) == "" && values[key] == "" {
			usesApplicationConfig = true
			break
		}
	}
	if usesApplicationConfig {
		environment := strings.ToLower(strings.TrimSpace(value("GOREST_TEST_MYSQL_APP_ENV", "GOREST_APP_ENV")))
		if environment != "development" && environment != "test" && environment != "testing" {
			return mysqlTestConfig{}, errors.New("GOREST_DB_* fallback requires GOREST_APP_ENV=development, test, or testing")
		}
		if strings.ToLower(strings.TrimSpace(value("GOREST_TEST_MYSQL_DRIVER", "GOREST_DB_DRIVER"))) != "mysql" {
			return mysqlTestConfig{}, errors.New("GOREST_DB_* fallback requires GOREST_DB_DRIVER=mysql")
		}
		if strings.EqualFold(strings.TrimSpace(value("GOREST_TEST_MYSQL_DISABLE", "GOREST_DB_DISABLE")), "true") {
			return mysqlTestConfig{}, errors.New("GOREST_DB_* fallback requires GOREST_DB_DISABLE=false")
		}
	}

	testValue := func(testKey, appKey string) string {
		if current := os.Getenv(testKey); current != "" {
			return current
		}
		if configured := values[testKey]; configured != "" {
			return configured
		}
		if usesApplicationConfig {
			if current := os.Getenv(appKey); current != "" {
				return current
			}
			return values[appKey]
		}
		return ""
	}

	config := mysqlTestConfig{
		host:     strings.TrimSpace(value("GOREST_TEST_MYSQL_HOST", "GOREST_DB_HOST")),
		port:     strings.TrimSpace(value("GOREST_TEST_MYSQL_PORT", "GOREST_DB_PORT")),
		database: strings.TrimSpace(value("GOREST_TEST_MYSQL_DATABASE", "GOREST_DB_NAME")),
		username: value("GOREST_TEST_MYSQL_USER", "GOREST_DB_USER"),
		password: value("GOREST_TEST_MYSQL_PASSWORD", "GOREST_DB_PASSWORD"),
		tls:      strings.TrimSpace(value("GOREST_TEST_MYSQL_TLS", "GOREST_DB_TLS")),
		caPath:   testValue("GOREST_TEST_MYSQL_CA_FILE", "GOREST_DB_CA_PATH"),
		certPath: testValue("GOREST_TEST_MYSQL_CERT_FILE", "GOREST_DB_CERT_PATH"),
		keyPath:  testValue("GOREST_TEST_MYSQL_KEY_FILE", "GOREST_DB_KEY_PATH"),
	}
	for _, path := range []*string{&config.caPath, &config.certPath, &config.keyPath} {
		if *path != "" {
			resolved, err := resolveMySQLTestPath(*path)
			if err != nil {
				return mysqlTestConfig{}, fmt.Errorf("resolve MySQL integration TLS file: %w", err)
			}
			*path = resolved
		}
	}
	if config.port == "" {
		config.port = "3306"
	}
	if config.tls == "" {
		config.tls = "disable"
	}

	for key, value := range map[string]string{
		"GOREST_TEST_MYSQL_HOST":     config.host,
		"GOREST_TEST_MYSQL_DATABASE": config.database,
		"GOREST_TEST_MYSQL_USER":     config.username,
		"GOREST_TEST_MYSQL_PASSWORD": config.password,
	} {
		if value == "" {
			return mysqlTestConfig{}, fmt.Errorf("%s is required", key)
		}
	}
	if config.tls != "true" && config.tls != "disable" {
		return mysqlTestConfig{}, errors.New("GOREST_TEST_MYSQL_TLS must be either true or disable")
	}
	if !strings.HasSuffix(strings.ToLower(config.database), "_test") &&
		!strings.HasSuffix(strings.ToLower(config.database), "_testing") &&
		value("GOREST_TEST_MYSQL_ALLOW_NON_TEST_DATABASE", "GOREST_TEST_MYSQL_ALLOW_NON_TEST_DATABASE") != config.database {
		return mysqlTestConfig{}, errors.New("integration database name must end in _test or _testing, or be explicitly allowed with GOREST_TEST_MYSQL_ALLOW_NON_TEST_DATABASE")
	}
	if value("GOREST_TEST_MYSQL_CONFIRM_DATABASE", "GOREST_TEST_MYSQL_CONFIRM_DATABASE") != config.database {
		return mysqlTestConfig{}, errors.New("set GOREST_TEST_MYSQL_CONFIRM_DATABASE to the exact test database name to allow test data writes")
	}

	return config, nil
}

func (config mysqlTestConfig) appEnvironment(serverPort, logDirectory, signatureKey string) map[string]string {
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
		"GOREST_DB_DRIVER":               "mysql",
		"GOREST_DB_HOST":                 config.host,
		"GOREST_DB_PORT":                 config.port,
		"GOREST_DB_NAME":                 config.database,
		"GOREST_DB_USER":                 config.username,
		"GOREST_DB_PASSWORD":             config.password,
		"GOREST_DB_TLS":                  config.tls,
		"GOREST_DB_CA_PATH":              config.caPath,
		"GOREST_DB_CERT_PATH":            config.certPath,
		"GOREST_DB_KEY_PATH":             config.keyPath,
		"GOREST_DB_MAX_IDLE_CONNECTIONS": "0",
		"GOREST_DB_MIGRATION":            "false",
		"GOREST_DB_SEEDER":               "false",
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
