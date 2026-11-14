//go:build integration

package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/elsyahtech/gorest/database"
)

func configureScyllaTestEnvironment(t *testing.T) {
	t.Helper()
	for key, value := range map[string]string{
		"GOREST_APP_ENV":                             "testing",
		"GOREST_TEST_SCYLLA_HOST":                    "scylla.example.test",
		"GOREST_TEST_SCYLLA_PORT":                    "",
		"GOREST_TEST_SCYLLA_KEYSPACE":                "gorest_starter_test",
		"GOREST_TEST_SCYLLA_USER":                    "gorest_test",
		"GOREST_TEST_SCYLLA_PASSWORD":                "integration-password",
		"GOREST_TEST_SCYLLA_TLS":                     "",
		"GOREST_TEST_SCYLLA_CA_FILE":                 "",
		"GOREST_TEST_SCYLLA_CERT_FILE":               "",
		"GOREST_TEST_SCYLLA_KEY_FILE":                "",
		"GOREST_TEST_SCYLLA_CONFIRM_KEYSPACE":        "gorest_starter_test",
		"GOREST_TEST_SCYLLA_ALLOW_NON_TEST_KEYSPACE": "",
		"GOREST_DB_DRIVER":                           "",
		"GOREST_DB_DISABLE":                          "",
		"GOREST_DB_HOST":                             "",
		"GOREST_DB_PORT":                             "",
		"GOREST_DB_NAME":                             "",
		"GOREST_DB_USER":                             "",
		"GOREST_DB_PASSWORD":                         "",
		"GOREST_DB_TLS":                              "",
		"GOREST_DB_CA_PATH":                          "",
		"GOREST_DB_CERT_PATH":                        "",
		"GOREST_DB_KEY_PATH":                         "",
	} {
		t.Setenv(key, value)
	}
}

func TestLoadScyllaTestConfigAcceptsDedicatedKeyspace(t *testing.T) {
	configureScyllaTestEnvironment(t)

	config, err := loadScyllaTestConfig()
	if err != nil {
		t.Fatalf("loadScyllaTestConfig() error = %v", err)
	}
	if config.host != "scylla.example.test" || config.port != database.SCYLLA_PORT ||
		config.keyspace != "gorest_starter_test" || config.username != "gorest_test" ||
		config.tls != database.TLSDisabled {
		t.Fatalf("ScyllaDB test config = (%+v), want configured keyspace and defaults", config)
	}
	environment := config.appEnvironment("3000", t.TempDir(), "test-signature-key")
	if environment["GOREST_DB_DRIVER"] != database.SCYLLA ||
		environment["GOREST_DB_NAME"] != config.keyspace ||
		environment["GOREST_DB_MIGRATION"] != "false" ||
		environment["GOREST_DB_SEEDER"] != "false" ||
		environment["GOREST_REDIS_DISABLE"] != "true" {
		t.Fatal("application environment does not select ScyllaDB and disable implicit startup writes")
	}
}

func TestLoadScyllaTestConfigRequiresExactKeyspaceConfirmation(t *testing.T) {
	configureScyllaTestEnvironment(t)
	t.Setenv("GOREST_TEST_SCYLLA_CONFIRM_KEYSPACE", "another_keyspace_test")

	if _, err := loadScyllaTestConfig(); err == nil {
		t.Fatal("loadScyllaTestConfig() accepted a keyspace without exact confirmation")
	}
}

func TestLoadScyllaTestConfigRequiresNonTestKeyspaceOptIn(t *testing.T) {
	configureScyllaTestEnvironment(t)
	t.Setenv("GOREST_TEST_SCYLLA_KEYSPACE", "gorest_starter")
	t.Setenv("GOREST_TEST_SCYLLA_CONFIRM_KEYSPACE", "gorest_starter")

	if _, err := loadScyllaTestConfig(); err == nil {
		t.Fatal("loadScyllaTestConfig() accepted an un-suffixed keyspace without explicit opt-in")
	}

	t.Setenv("GOREST_TEST_SCYLLA_ALLOW_NON_TEST_KEYSPACE", "gorest_starter")
	if _, err := loadScyllaTestConfig(); err != nil {
		t.Fatalf("loadScyllaTestConfig() with exact non-test keyspace opt-in error = %v", err)
	}
}

func TestLoadScyllaTestConfigRejectsProduction(t *testing.T) {
	configureScyllaTestEnvironment(t)
	t.Setenv("GOREST_APP_ENV", "production")

	if _, err := loadScyllaTestConfig(); err == nil {
		t.Fatal("loadScyllaTestConfig() accepted production mode")
	}
}

func TestLoadScyllaTestConfigRequiresPairedCredentialsAndCertificate(t *testing.T) {
	configureScyllaTestEnvironment(t)
	t.Setenv("GOREST_TEST_SCYLLA_PASSWORD", "")

	if _, err := loadScyllaTestConfig(); err == nil {
		t.Fatal("loadScyllaTestConfig() accepted a username without a password")
	}

	configureScyllaTestEnvironment(t)
	t.Setenv("GOREST_TEST_SCYLLA_TLS", "enable")
	t.Setenv("GOREST_TEST_SCYLLA_CERT_FILE", "/missing/client.crt")
	if _, err := loadScyllaTestConfig(); err == nil {
		t.Fatal("loadScyllaTestConfig() accepted an unpaired client certificate")
	}
}

func TestLoadScyllaTestConfigRequiresSafeTLSMode(t *testing.T) {
	configureScyllaTestEnvironment(t)
	t.Setenv("GOREST_TEST_SCYLLA_TLS", "true")

	if _, err := loadScyllaTestConfig(); err == nil {
		t.Fatal("loadScyllaTestConfig() accepted an unsupported TLS mode")
	}

	configureScyllaTestEnvironment(t)
	t.Setenv("GOREST_TEST_SCYLLA_TLS", "enable")
	t.Setenv("GOREST_TEST_SCYLLA_CA_FILE", "/missing/ca.pem")
	if _, err := loadScyllaTestConfig(); err == nil {
		t.Fatal("loadScyllaTestConfig() accepted an inaccessible CA file")
	}
}

func TestResolveScyllaTestPathRequiresRegularFile(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "ca.pem")
	if err := os.WriteFile(path, []byte("certificate"), 0o600); err != nil {
		t.Fatalf("write temporary CA file: %v", err)
	}

	resolved, err := resolveScyllaTestPath(root, path)
	if err != nil || resolved != path {
		t.Fatalf("resolveScyllaTestPath() = (%q, %v), want %q", resolved, err, path)
	}
	if _, err := resolveScyllaTestPath(root, root); err == nil {
		t.Fatal("resolveScyllaTestPath() accepted a directory")
	}
}

func TestLoadScyllaTestConfigRequiresRuntimeFallbackOptIn(t *testing.T) {
	configureScyllaTestEnvironment(t)
	for _, key := range []string{
		"GOREST_TEST_SCYLLA_HOST",
		"GOREST_TEST_SCYLLA_KEYSPACE",
		"GOREST_TEST_SCYLLA_USER",
		"GOREST_TEST_SCYLLA_PASSWORD",
	} {
		t.Setenv(key, "")
	}
	t.Setenv("GOREST_DB_DRIVER", "scylla")
	t.Setenv("GOREST_DB_DISABLE", "false")
	t.Setenv("GOREST_DB_HOST", "scylla.example.test")
	t.Setenv("GOREST_DB_PORT", "9142")
	t.Setenv("GOREST_DB_NAME", "gorest_starter_test")
	t.Setenv("GOREST_DB_USER", "gorest_test")
	t.Setenv("GOREST_DB_PASSWORD", "integration-password")
	t.Setenv("GOREST_DB_TLS", "enable")
	t.Setenv("GOREST_TEST_SCYLLA_CONFIRM_KEYSPACE", "gorest_starter_test")

	config, err := loadScyllaTestConfig()
	if err != nil {
		t.Fatalf("loadScyllaTestConfig() with safe runtime fallback error = %v", err)
	}
	if config.host != "scylla.example.test" || config.port != "9142" ||
		config.keyspace != "gorest_starter_test" || config.tls != database.TLSEnabled {
		t.Fatalf("runtime fallback config = (%+v), want application endpoint", config)
	}

	t.Setenv("GOREST_APP_ENV", "production")
	if _, err := loadScyllaTestConfig(); err == nil {
		t.Fatal("loadScyllaTestConfig() accepted runtime fallback in production")
	}
}
