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

type oracleTestConfig struct {
	host              string
	port              string
	serviceName       string
	username          string
	password          string
	tls               string
	walletPath        string
	libDir            string
	databaseDirectory string
}

func loadOracleTestConfig() (oracleTestConfig, error) {
	values, err := godotenv.Read("../.env")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return oracleTestConfig{}, fmt.Errorf("read .env: %w", err)
	}
	value := func(key string) string {
		if current, ok := os.LookupEnv(key); ok {
			return current
		}
		return values[key]
	}

	if strings.EqualFold(strings.TrimSpace(value("GOREST_APP_ENV")), "production") {
		return oracleTestConfig{}, errors.New("Oracle integration tests refuse to run with GOREST_APP_ENV=production")
	}
	if !strings.EqualFold(strings.TrimSpace(value("GOREST_DB_DRIVER")), database.ORACLE) {
		return oracleTestConfig{}, errors.New("GOREST_DB_DRIVER=oracle is required for Oracle integration tests")
	}
	if strings.EqualFold(strings.TrimSpace(value("GOREST_DB_DISABLE")), "true") {
		return oracleTestConfig{}, errors.New("GOREST_DB_DISABLE=false is required for Oracle integration tests")
	}

	config := oracleTestConfig{
		host:        strings.TrimSpace(value("GOREST_DB_HOST")),
		port:        strings.TrimSpace(value("GOREST_DB_PORT")),
		serviceName: strings.TrimSpace(value("GOREST_DB_NAME")),
		username:    strings.TrimSpace(value("GOREST_DB_USER")),
		password:    strings.TrimSpace(value("GOREST_DB_PASSWORD")),
		tls:         strings.ToLower(strings.TrimSpace(value("GOREST_DB_TLS"))),
		walletPath:  strings.TrimSpace(value("GOREST_DB_CA_PATH")),
		libDir:      strings.TrimSpace(value("GOREST_DB_LIB_DIR_PATH")),
	}
	if config.host == "" {
		config.host = database.DefaultDatabaseHost
	}
	if config.port == "" {
		config.port = database.ORACLE_PORT
	}
	if config.tls == "" {
		config.tls = database.TLSDisabled
	}
	if config.tls != database.TLSEnabled && config.tls != database.TLSDisabled {
		return oracleTestConfig{}, errors.New("GOREST_DB_TLS must be enable or disable for Oracle integration tests")
	}
	for key, configured := range map[string]string{
		"GOREST_DB_NAME (Oracle service name)": config.serviceName,
		"GOREST_DB_USER":                       config.username,
		"GOREST_DB_PASSWORD":                   config.password,
	} {
		if configured == "" {
			return oracleTestConfig{}, fmt.Errorf("%s is required for Oracle integration tests", key)
		}
	}
	switch strings.ToUpper(config.username) {
	case "SYS", "SYSTEM", "DBSNMP", "OUTLN", "SYSBACKUP", "SYSDG", "SYSKM", "SYSRAC":
		return oracleTestConfig{}, errors.New("Oracle integration tests refuse to run as a built-in administrative account")
	}
	if value("GOREST_TEST_ORACLE_CONFIRM_USER") != config.username {
		return oracleTestConfig{}, errors.New("set GOREST_TEST_ORACLE_CONFIRM_USER to exactly match the dedicated non-production Oracle schema username before running schema-changing tests")
	}
	if config.tls == database.TLSEnabled && config.walletPath == "" {
		return oracleTestConfig{}, errors.New("GOREST_DB_CA_PATH must point to an Oracle wallet directory when GOREST_DB_TLS=enable")
	}

	root, err := filepath.Abs("..")
	if err != nil {
		return oracleTestConfig{}, fmt.Errorf("resolve repository root for Oracle paths: %w", err)
	}
	config.walletPath, err = resolveOraclePath(root, config.walletPath)
	if err != nil {
		return oracleTestConfig{}, fmt.Errorf("resolve Oracle wallet directory: %w", err)
	}
	config.libDir, err = resolveOraclePath(root, config.libDir)
	if err != nil {
		return oracleTestConfig{}, fmt.Errorf("resolve Oracle Instant Client directory: %w", err)
	}
	if config.tls == database.TLSEnabled {
		if err := requireDirectory(config.walletPath); err != nil {
			return oracleTestConfig{}, fmt.Errorf("access configured Oracle wallet directory: %w", err)
		}
	}
	if config.libDir != "" {
		if err := requireDirectory(config.libDir); err != nil {
			return oracleTestConfig{}, fmt.Errorf("access configured Oracle Instant Client directory: %w", err)
		}
	}
	config.databaseDirectory = filepath.Join(root, "app", "database")
	return config, nil
}

func resolveOraclePath(root, path string) (string, error) {
	if path == "" || filepath.IsAbs(path) {
		return path, nil
	}
	return filepath.Abs(filepath.Join(root, path))
}

func requireDirectory(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%q is not a directory", path)
	}
	return nil
}

func (config oracleTestConfig) appEnvironment(serverPort, logDirectory, signatureKey string) map[string]string {
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
		"GOREST_DB_DRIVER":                    database.ORACLE,
		"GOREST_DB_HOST":                      config.host,
		"GOREST_DB_PORT":                      config.port,
		"GOREST_DB_NAME":                      config.serviceName,
		"GOREST_DB_USER":                      config.username,
		"GOREST_DB_PASSWORD":                  config.password,
		"GOREST_DB_TLS":                       config.tls,
		"GOREST_DB_CA_PATH":                   config.walletPath,
		"GOREST_DB_LIB_DIR_PATH":              config.libDir,
		"GOREST_DB_DIRECTORY":                 config.databaseDirectory,
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
