//go:build integration

package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/elsyahtech/gorest/database"
)

func configureMongoTestEnvironment(t *testing.T) {
	t.Helper()
	t.Setenv("GOREST_APP_ENV", "testing")
	t.Setenv("GOREST_TEST_MONGO_HOST", "mongo.example.test")
	t.Setenv("GOREST_TEST_MONGO_PORT", "")
	t.Setenv("GOREST_TEST_MONGO_DATABASE", "gorest_starter_test")
	t.Setenv("GOREST_TEST_MONGO_USER", "gorest_test")
	t.Setenv("GOREST_TEST_MONGO_PASSWORD", "integration-password")
	t.Setenv("GOREST_TEST_MONGO_AUTH_SOURCE", "")
	t.Setenv("GOREST_TEST_MONGO_TLS", "")
	t.Setenv("GOREST_TEST_MONGO_CA_FILE", "")
	t.Setenv("GOREST_TEST_MONGO_CERT_FILE", "")
	t.Setenv("GOREST_TEST_MONGO_KEY_FILE", "")
	t.Setenv("GOREST_TEST_MONGO_CONFIRM_DATABASE", "gorest_starter_test")
	t.Setenv("GOREST_TEST_MONGO_ALLOW_NON_TEST_DATABASE", "")
}

func TestLoadMongoTestConfigAcceptsDedicatedDatabase(t *testing.T) {
	configureMongoTestEnvironment(t)

	config, err := loadMongoTestConfig()
	if err != nil {
		t.Fatalf("loadMongoTestConfig() error = %v", err)
	}
	if config.host != "mongo.example.test" || config.port != database.MONGO_PORT ||
		config.database != "gorest_starter_test" || config.username != "gorest_test" ||
		config.authSource != "admin" || config.tls != database.TLSFalse {
		t.Fatalf("MongoDB test config = (%+v), want configured test database and defaults", config)
	}
	environment := config.appEnvironment("3000", t.TempDir(), "test-signature-key")
	if environment["GOREST_DB_DRIVER"] != database.MONGO ||
		environment["GOREST_DB_NAME"] != config.database ||
		environment["GOREST_DB_AUTH_SOURCE"] != "admin" ||
		environment["GOREST_DB_MIGRATION"] != "false" ||
		environment["GOREST_REDIS_DISABLE"] != "true" {
		t.Fatal("application environment does not select MongoDB and disable implicit startup writes")
	}
}

func TestLoadMongoTestConfigRequiresExactDatabaseConfirmation(t *testing.T) {
	configureMongoTestEnvironment(t)
	t.Setenv("GOREST_TEST_MONGO_CONFIRM_DATABASE", "another_database_test")

	if _, err := loadMongoTestConfig(); err == nil {
		t.Fatal("loadMongoTestConfig() accepted a database without exact confirmation")
	}
}

func TestLoadMongoTestConfigRequiresNonTestDatabaseOptIn(t *testing.T) {
	configureMongoTestEnvironment(t)
	t.Setenv("GOREST_TEST_MONGO_DATABASE", "gorest_starter")
	t.Setenv("GOREST_TEST_MONGO_CONFIRM_DATABASE", "gorest_starter")

	if _, err := loadMongoTestConfig(); err == nil {
		t.Fatal("loadMongoTestConfig() accepted an un-suffixed database without explicit opt-in")
	}

	t.Setenv("GOREST_TEST_MONGO_ALLOW_NON_TEST_DATABASE", "gorest_starter")
	if _, err := loadMongoTestConfig(); err != nil {
		t.Fatalf("loadMongoTestConfig() with exact non-test database opt-in error = %v", err)
	}
}

func TestLoadMongoTestConfigRejectsProduction(t *testing.T) {
	configureMongoTestEnvironment(t)
	t.Setenv("GOREST_APP_ENV", "production")

	if _, err := loadMongoTestConfig(); err == nil {
		t.Fatal("loadMongoTestConfig() accepted production mode")
	}
}

func TestLoadMongoTestConfigRequiresPairedClientCertificate(t *testing.T) {
	configureMongoTestEnvironment(t)
	t.Setenv("GOREST_TEST_MONGO_TLS", "true")
	t.Setenv("GOREST_TEST_MONGO_CERT_FILE", "/missing/client.crt")

	if _, err := loadMongoTestConfig(); err == nil {
		t.Fatal("loadMongoTestConfig() accepted an unpaired client certificate")
	}
}

func TestResolveMongoTestPathRequiresRegularFile(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "ca.pem")
	if err := os.WriteFile(path, []byte("certificate"), 0o600); err != nil {
		t.Fatalf("write temporary CA file: %v", err)
	}

	resolved, err := resolveMongoTestPath(root, path)
	if err != nil || resolved != path {
		t.Fatalf("resolveMongoTestPath() = (%q, %v), want %q", resolved, err, path)
	}
	if _, err := resolveMongoTestPath(root, root); err == nil {
		t.Fatal("resolveMongoTestPath() accepted a directory")
	}
}
