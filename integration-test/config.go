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

type mongoTestConfig struct {
	host        string
	port        string
	database    string
	username    string
	password    string
	authSource  string
	tls         string
	caPath      string
	certPath    string
	keyPath     string
	databaseDir string
}

func loadMongoTestConfig() (mongoTestConfig, error) {
	values, err := godotenv.Read("../.env")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return mongoTestConfig{}, fmt.Errorf("read .env: %w", err)
	}
	value := func(key string) string {
		if current, ok := os.LookupEnv(key); ok && current != "" {
			return current
		}
		return values[key]
	}

	if strings.EqualFold(strings.TrimSpace(value("GOREST_APP_ENV")), "production") {
		return mongoTestConfig{}, errors.New("MongoDB integration tests refuse to run with GOREST_APP_ENV=production")
	}

	testConnectionKeys := []string{
		"GOREST_TEST_MONGO_HOST",
		"GOREST_TEST_MONGO_DATABASE",
		"GOREST_TEST_MONGO_USER",
		"GOREST_TEST_MONGO_PASSWORD",
	}
	useRuntimeConfig := true
	for _, key := range testConnectionKeys {
		if value(key) != "" {
			useRuntimeConfig = false
			break
		}
	}

	runtimeValue := func(testKey, runtimeKey string) string {
		if configured := value(testKey); configured != "" {
			return configured
		}
		if useRuntimeConfig {
			return value(runtimeKey)
		}
		return ""
	}

	if useRuntimeConfig {
		environment := strings.ToLower(strings.TrimSpace(value("GOREST_APP_ENV")))
		if environment != "development" && environment != "test" && environment != "testing" {
			return mongoTestConfig{}, errors.New("GOREST_DB_* fallback requires GOREST_APP_ENV=development, test, or testing")
		}
		if !strings.EqualFold(strings.TrimSpace(value("GOREST_DB_DRIVER")), database.MONGO) {
			return mongoTestConfig{}, errors.New("GOREST_DB_* fallback requires GOREST_DB_DRIVER=mongo")
		}
		if strings.EqualFold(strings.TrimSpace(value("GOREST_DB_DISABLE")), "true") {
			return mongoTestConfig{}, errors.New("GOREST_DB_* fallback requires GOREST_DB_DISABLE=false")
		}
	}

	config := mongoTestConfig{
		host:       strings.TrimSpace(runtimeValue("GOREST_TEST_MONGO_HOST", "GOREST_DB_HOST")),
		port:       strings.TrimSpace(runtimeValue("GOREST_TEST_MONGO_PORT", "GOREST_DB_PORT")),
		database:   strings.TrimSpace(runtimeValue("GOREST_TEST_MONGO_DATABASE", "GOREST_DB_NAME")),
		username:   strings.TrimSpace(runtimeValue("GOREST_TEST_MONGO_USER", "GOREST_DB_USER")),
		password:   strings.TrimSpace(runtimeValue("GOREST_TEST_MONGO_PASSWORD", "GOREST_DB_PASSWORD")),
		authSource: strings.TrimSpace(runtimeValue("GOREST_TEST_MONGO_AUTH_SOURCE", "GOREST_DB_AUTH_SOURCE")),
		tls:        strings.ToLower(strings.TrimSpace(runtimeValue("GOREST_TEST_MONGO_TLS", "GOREST_DB_TLS"))),
		caPath:     strings.TrimSpace(runtimeValue("GOREST_TEST_MONGO_CA_FILE", "GOREST_DB_CA_PATH")),
		certPath:   strings.TrimSpace(runtimeValue("GOREST_TEST_MONGO_CERT_FILE", "GOREST_DB_CERT_PATH")),
		keyPath:    strings.TrimSpace(runtimeValue("GOREST_TEST_MONGO_KEY_FILE", "GOREST_DB_KEY_PATH")),
	}
	if config.host == "" {
		config.host = database.DefaultDatabaseHost
	}
	if config.port == "" {
		config.port = database.MONGO_PORT
	}
	if config.authSource == "" {
		config.authSource = "admin"
	}
	if config.tls == "" {
		config.tls = database.TLSFalse
	}
	if config.tls != database.TLSFalse && config.tls != "true" {
		return mongoTestConfig{}, errors.New("MongoDB TLS must be true or false")
	}
	if config.database == "" {
		return mongoTestConfig{}, errors.New("GOREST_TEST_MONGO_DATABASE (or GOREST_DB_NAME) is required")
	}
	if (config.username == "") != (config.password == "") {
		return mongoTestConfig{}, errors.New("MongoDB test username and password must be configured together")
	}
	if (config.certPath == "") != (config.keyPath == "") {
		return mongoTestConfig{}, errors.New("MongoDB test client certificate and key must be configured together")
	}
	if config.certPath != "" && config.tls != "true" {
		return mongoTestConfig{}, errors.New("MongoDB TLS=true is required when client certificates are configured")
	}

	confirmation := strings.TrimSpace(value("GOREST_TEST_MONGO_CONFIRM_DATABASE"))
	if confirmation != config.database {
		return mongoTestConfig{}, errors.New("set GOREST_TEST_MONGO_CONFIRM_DATABASE to exactly match the dedicated non-production MongoDB test database before running database-writing tests")
	}
	if !strings.HasSuffix(strings.ToLower(config.database), "_test") &&
		!strings.HasSuffix(strings.ToLower(config.database), "_testing") &&
		strings.TrimSpace(value("GOREST_TEST_MONGO_ALLOW_NON_TEST_DATABASE")) != config.database {
		return mongoTestConfig{}, errors.New("MongoDB integration database must end in _test or _testing, or be explicitly allowed with GOREST_TEST_MONGO_ALLOW_NON_TEST_DATABASE")
	}

	root, err := filepath.Abs("..")
	if err != nil {
		return mongoTestConfig{}, fmt.Errorf("resolve repository root for MongoDB paths: %w", err)
	}
	for _, configuredPath := range []*string{&config.caPath, &config.certPath, &config.keyPath} {
		if *configuredPath == "" {
			continue
		}
		*configuredPath, err = resolveMongoTestPath(root, *configuredPath)
		if err != nil {
			return mongoTestConfig{}, fmt.Errorf("resolve MongoDB TLS file: %w", err)
		}
	}
	config.databaseDir = filepath.Join(root, "app", "database")
	return config, nil
}

func resolveMongoTestPath(root, path string) (string, error) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(root, path)
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(absolute)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("%q is not a regular file", absolute)
	}
	return absolute, nil
}

func (config mongoTestConfig) appEnvironment(serverPort, logDirectory, signatureKey string) map[string]string {
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
