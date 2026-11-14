//go:build integration

package integration_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/config"
	"github.com/elsyahtech/gorest/database"
	"github.com/elsyahtech/gorest/log"
	"github.com/joho/godotenv"
)

type scyllaTestConfig struct {
	host        string
	port        string
	keyspace    string
	username    string
	password    string
	tls         string
	caPath      string
	certPath    string
	keyPath     string
	databaseDir string
}

func loadScyllaTestConfig() (scyllaTestConfig, error) {
	values, err := godotenv.Read("../.env")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return scyllaTestConfig{}, fmt.Errorf("read .env: %w", err)
	}
	value := func(key string) string {
		if current, ok := os.LookupEnv(key); ok {
			return current
		}
		return values[key]
	}

	if strings.EqualFold(strings.TrimSpace(value("GOREST_APP_ENV")), gorest.Production) {
		return scyllaTestConfig{}, errors.New("ScyllaDB integration tests refuse to run with GOREST_APP_ENV=production")
	}

	useRuntimeConfig := true
	for _, key := range []string{
		"GOREST_TEST_SCYLLA_HOST",
		"GOREST_TEST_SCYLLA_KEYSPACE",
		"GOREST_TEST_SCYLLA_USER",
		"GOREST_TEST_SCYLLA_PASSWORD",
	} {
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
			return scyllaTestConfig{}, errors.New("GOREST_DB_* fallback requires GOREST_APP_ENV=development, test, or testing")
		}
		if !strings.EqualFold(strings.TrimSpace(value("GOREST_DB_DRIVER")), database.SCYLLA) {
			return scyllaTestConfig{}, errors.New("GOREST_DB_* fallback requires GOREST_DB_DRIVER=scylla")
		}
		if strings.EqualFold(strings.TrimSpace(value("GOREST_DB_DISABLE")), "true") {
			return scyllaTestConfig{}, errors.New("GOREST_DB_* fallback requires GOREST_DB_DISABLE=false")
		}
	}

	config := scyllaTestConfig{
		host:     strings.TrimSpace(runtimeValue("GOREST_TEST_SCYLLA_HOST", "GOREST_DB_HOST")),
		port:     strings.TrimSpace(runtimeValue("GOREST_TEST_SCYLLA_PORT", "GOREST_DB_PORT")),
		keyspace: strings.TrimSpace(runtimeValue("GOREST_TEST_SCYLLA_KEYSPACE", "GOREST_DB_NAME")),
		username: strings.TrimSpace(runtimeValue("GOREST_TEST_SCYLLA_USER", "GOREST_DB_USER")),
		password: strings.TrimSpace(runtimeValue("GOREST_TEST_SCYLLA_PASSWORD", "GOREST_DB_PASSWORD")),
		tls:      strings.ToLower(strings.TrimSpace(runtimeValue("GOREST_TEST_SCYLLA_TLS", "GOREST_DB_TLS"))),
		caPath:   strings.TrimSpace(runtimeValue("GOREST_TEST_SCYLLA_CA_FILE", "GOREST_DB_CA_PATH")),
		certPath: strings.TrimSpace(runtimeValue("GOREST_TEST_SCYLLA_CERT_FILE", "GOREST_DB_CERT_PATH")),
		keyPath:  strings.TrimSpace(runtimeValue("GOREST_TEST_SCYLLA_KEY_FILE", "GOREST_DB_KEY_PATH")),
	}
	if config.host == "" {
		config.host = database.DefaultDatabaseHost
	}
	if config.port == "" {
		config.port = database.SCYLLA_PORT
	}
	if config.tls == "" {
		config.tls = database.TLSDisabled
	}
	if config.keyspace == "" {
		return scyllaTestConfig{}, errors.New("GOREST_TEST_SCYLLA_KEYSPACE (or GOREST_DB_NAME) is required")
	}
	if config.tls != database.TLSDisabled && config.tls != database.TLSEnabled {
		return scyllaTestConfig{}, errors.New("ScyllaDB TLS must be enable or disable")
	}
	if (config.username == "") != (config.password == "") {
		return scyllaTestConfig{}, errors.New("ScyllaDB test username and password must be configured together")
	}
	if (config.certPath == "") != (config.keyPath == "") {
		return scyllaTestConfig{}, errors.New("ScyllaDB test client certificate and key must be configured together")
	}
	if (config.caPath != "" || config.certPath != "") && config.tls != database.TLSEnabled {
		return scyllaTestConfig{}, errors.New("ScyllaDB TLS=enable is required when TLS files are configured")
	}

	confirmation := strings.TrimSpace(value("GOREST_TEST_SCYLLA_CONFIRM_KEYSPACE"))
	if confirmation != config.keyspace {
		return scyllaTestConfig{}, errors.New("set GOREST_TEST_SCYLLA_CONFIRM_KEYSPACE to exactly match the dedicated non-production ScyllaDB test keyspace before running database-writing tests")
	}
	if !strings.HasSuffix(strings.ToLower(config.keyspace), "_test") &&
		!strings.HasSuffix(strings.ToLower(config.keyspace), "_testing") &&
		strings.TrimSpace(value("GOREST_TEST_SCYLLA_ALLOW_NON_TEST_KEYSPACE")) != config.keyspace {
		return scyllaTestConfig{}, errors.New("ScyllaDB integration keyspace must end in _test or _testing, or be explicitly allowed with GOREST_TEST_SCYLLA_ALLOW_NON_TEST_KEYSPACE")
	}

	root, err := filepath.Abs("..")
	if err != nil {
		return scyllaTestConfig{}, fmt.Errorf("resolve repository root for ScyllaDB paths: %w", err)
	}
	for _, configuredPath := range []*string{&config.caPath, &config.certPath, &config.keyPath} {
		if *configuredPath == "" {
			continue
		}
		*configuredPath, err = resolveScyllaTestPath(root, *configuredPath)
		if err != nil {
			return scyllaTestConfig{}, fmt.Errorf("resolve ScyllaDB TLS file: %w", err)
		}
	}
	config.databaseDir = filepath.Join(root, "app", "database")
	return config, nil
}

func resolveScyllaTestPath(root, path string) (string, error) {
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

func applyScyllaTestConfig(t *testing.T, config scyllaTestConfig) {
	t.Helper()
	for key, value := range map[string]string{
		"GOREST_APP_ENV":                      "testing",
		"GOREST_DB_DISABLE":                   "false",
		"GOREST_DB_DRIVER":                    database.SCYLLA,
		"GOREST_DB_HOST":                      config.host,
		"GOREST_DB_PORT":                      config.port,
		"GOREST_DB_NAME":                      config.keyspace,
		"GOREST_DB_USER":                      config.username,
		"GOREST_DB_PASSWORD":                  config.password,
		"GOREST_DB_TLS":                       config.tls,
		"GOREST_DB_CA_PATH":                   config.caPath,
		"GOREST_DB_CERT_PATH":                 config.certPath,
		"GOREST_DB_KEY_PATH":                  config.keyPath,
		"GOREST_DB_DIRECTORY":                 config.databaseDir,
		"GOREST_DB_MIGRATION":                 "false",
		"GOREST_DB_SEEDER":                    "false",
		"GOREST_DB_ALLOW_INSECURE_PRODUCTION": "false",
	} {
		t.Setenv(key, value)
	}
}

func prepareScyllaIntegrationSchema(t *testing.T, testConfig scyllaTestConfig) {
	t.Helper()
	applyScyllaTestConfig(t, testConfig)

	dbConfig, err := config.DatabaseConfig()
	if err != nil {
		t.Fatalf("load ScyllaDB migration configuration: %v", err)
	}
	dbConfig.DatabaseDirectory = testConfig.databaseDir
	dbConfig.Migration = true
	dbConfig.Seeder = true

	app := gorest.New(gorest.ConfigDefault)
	app.Use(log.Run(log.ConfigDefault), database.Run(*dbConfig))
	t.Cleanup(func() {
		if _, _, err := app.Close(); err != nil {
			t.Errorf("close ScyllaDB migration application: %v", err)
		}
	})
	if err := app.RunDatabaseStartupTasks(); err != nil {
		t.Fatalf("apply ScyllaDB integration migration and seeder: %v", err)
	}
}

func (config scyllaTestConfig) appEnvironment(serverPort, logDirectory, signatureKey string) map[string]string {
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
		"GOREST_DB_DRIVER":                    database.SCYLLA,
		"GOREST_DB_HOST":                      config.host,
		"GOREST_DB_PORT":                      config.port,
		"GOREST_DB_NAME":                      config.keyspace,
		"GOREST_DB_USER":                      config.username,
		"GOREST_DB_PASSWORD":                  config.password,
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
