//go:build integration

package integration_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/elsyahtech/gorest/database"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func TestMySQLMutualTLSRequiresClientCertificate(t *testing.T) {
	values, err := godotenv.Read("../.env")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("read .env for MySQL mTLS test: %v", err)
	}
	value := func(key string) string {
		if current := os.Getenv(key); current != "" {
			return current
		}
		return values[key]
	}
	setting := func(testKey, runtimeKey string) string {
		if configured := value(testKey); configured != "" {
			return configured
		}
		return value(runtimeKey)
	}

	requiredSettings := map[string]string{
		"GOREST_TEST_MYSQL_MTLS_HOST":             "GOREST_DB_HOST",
		"GOREST_TEST_MYSQL_MTLS_DATABASE":         "GOREST_DB_NAME",
		"GOREST_TEST_MYSQL_MTLS_USER":             "GOREST_DB_USER",
		"GOREST_TEST_MYSQL_MTLS_PASSWORD":         "GOREST_DB_PASSWORD",
		"GOREST_TEST_MYSQL_MTLS_CA_FILE":          "GOREST_DB_CA_PATH",
		"GOREST_TEST_MYSQL_MTLS_CLIENT_CERT_FILE": "GOREST_DB_CERT_PATH",
		"GOREST_TEST_MYSQL_MTLS_CLIENT_KEY_FILE":  "GOREST_DB_KEY_PATH",
	}
	configuredTestSettings := 0
	for testKey := range requiredSettings {
		if value(testKey) != "" {
			configuredTestSettings++
		}
	}
	useRuntimeSettings := configuredTestSettings == 0
	if useRuntimeSettings {
		configuredRuntimeSettings := 0
		for _, runtimeKey := range requiredSettings {
			if value(runtimeKey) != "" {
				configuredRuntimeSettings++
			}
		}
		if configuredRuntimeSettings == 0 {
			t.Skip("configure GOREST_TEST_MYSQL_MTLS_* or GOREST_DB_* TLS certificate settings to run MySQL mTLS integration test")
		}
		if !strings.EqualFold(strings.TrimSpace(value("GOREST_APP_ENV")), "development") &&
			!strings.EqualFold(strings.TrimSpace(value("GOREST_APP_ENV")), "test") &&
			!strings.EqualFold(strings.TrimSpace(value("GOREST_APP_ENV")), "testing") {
			t.Fatal("using GOREST_DB_* for MySQL mTLS requires GOREST_APP_ENV=development, test, or testing")
		}
		if configuredRuntimeSettings != len(requiredSettings) {
			for _, runtimeKey := range requiredSettings {
				if value(runtimeKey) == "" {
					t.Errorf("%s is required when runtime MySQL mTLS settings are used", runtimeKey)
				}
			}
			t.FailNow()
		}
	}
	if !useRuntimeSettings && configuredTestSettings != len(requiredSettings) {
		for testKey := range requiredSettings {
			if value(testKey) == "" {
				t.Errorf("%s is required when dedicated MySQL mTLS settings are configured", testKey)
			}
		}
		t.FailNow()
	}

	if tlsMode := strings.TrimSpace(setting("GOREST_TEST_MYSQL_MTLS_TLS", "GOREST_DB_TLS")); tlsMode != "true" {
		t.Fatalf("MySQL mTLS requires TLS=true; configured mode is %q", tlsMode)
	}
	for testKey, runtimeKey := range requiredSettings {
		if setting(testKey, runtimeKey) == "" {
			t.Fatalf("%s is required for MySQL mTLS", testKey)
		}
	}
	resolvedPaths := make(map[string]string, 3)
	for _, settingKey := range []string{
		"GOREST_TEST_MYSQL_MTLS_CA_FILE",
		"GOREST_TEST_MYSQL_MTLS_CLIENT_CERT_FILE",
		"GOREST_TEST_MYSQL_MTLS_CLIENT_KEY_FILE",
	} {
		path, err := resolveMySQLTestPath(setting(settingKey, requiredSettings[settingKey]))
		if err != nil {
			t.Fatalf("resolve configured MySQL mTLS certificate file: %v", err)
		}
		resolvedPaths[settingKey] = path
	}

	connect := func(withClientCertificate bool) error {
		t.Helper()
		port := strings.TrimSpace(setting("GOREST_TEST_MYSQL_MTLS_PORT", "GOREST_DB_PORT"))
		if port == "" {
			port = "3306"
		}
		config := database.ConfigDefault
		config.Disable = false
		config.Driver = database.MYSQL
		config.Host = setting("GOREST_TEST_MYSQL_MTLS_HOST", "GOREST_DB_HOST")
		config.Port = port
		config.Name = setting("GOREST_TEST_MYSQL_MTLS_DATABASE", "GOREST_DB_NAME")
		config.Username = setting("GOREST_TEST_MYSQL_MTLS_USER", "GOREST_DB_USER")
		config.Password = setting("GOREST_TEST_MYSQL_MTLS_PASSWORD", "GOREST_DB_PASSWORD")
		config.TLS = database.TLSTrue
		config.SSLCAPath = resolvedPaths["GOREST_TEST_MYSQL_MTLS_CA_FILE"]
		if withClientCertificate {
			config.SSLCertPath = resolvedPaths["GOREST_TEST_MYSQL_MTLS_CLIENT_CERT_FILE"]
			config.SSLKeyPath = resolvedPaths["GOREST_TEST_MYSQL_MTLS_CLIENT_KEY_FILE"]
		}
		config.Timeout = 10 * time.Second

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		db, _, err := database.New(ctx, &config, "UTC")
		if err != nil {
			return err
		}
		defer func() {
			if err := db.SQL.Close(); err != nil {
				t.Errorf("close MySQL mTLS connection: %v", err)
			}
		}()
		if err := db.SQL.PingContext(ctx); err != nil {
			return err
		}

		if withClientCertificate {
			var status, cipher string
			if err := db.SQL.QueryRowContext(ctx, "SHOW SESSION STATUS LIKE 'Ssl_cipher'").Scan(&status, &cipher); err != nil {
				t.Fatalf("query MySQL mTLS status: %v", err)
			}
			if status != "Ssl_cipher" || cipher == "" {
				t.Fatalf("MySQL mTLS connection did not negotiate TLS (status=%q, cipher=%q)", status, cipher)
			}
		}

		return nil
	}

	err = connect(false)
	var mysqlErr *mysql.MySQLError
	if !errors.As(err, &mysqlErr) {
		t.Fatalf("MySQL connection without client certificate error = %v, want server access denial", err)
	}
	if mysqlErr.Number != 1045 {
		t.Fatalf("MySQL connection without client certificate error number = %d, want 1045 (access denied)", mysqlErr.Number)
	}
	if err := connect(true); err != nil {
		t.Fatalf("MySQL connection with client certificate failed: %v", err)
	}
}

func resolveMySQLTestPath(path string) (string, error) {
	if filepath.IsAbs(path) {
		if _, err := os.Stat(path); err != nil {
			return "", errors.New("configured MySQL TLS file is not accessible")
		}
		return path, nil
	}
	for _, candidate := range []string{path, filepath.Join("..", path)} {
		absolute, err := filepath.Abs(candidate)
		if err != nil {
			return "", errors.New("resolve configured MySQL TLS file path")
		}
		if _, err := os.Stat(absolute); err == nil {
			return absolute, nil
		}
	}
	return "", errors.New("configured MySQL TLS file is not accessible")
}
