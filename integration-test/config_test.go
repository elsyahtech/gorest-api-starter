//go:build integration

package integration_test

import "testing"

func TestLoadPostgreSQLTestConfigAllowsMutualTLS(t *testing.T) {
	t.Setenv("GOREST_APP_ENV", "testing")
	t.Setenv("GOREST_DB_DISABLE", "false")
	t.Setenv("GOREST_DB_DRIVER", "postgres")
	t.Setenv("GOREST_DB_HOST", "postgres.example.test")
	t.Setenv("GOREST_DB_PORT", "5432")
	t.Setenv("GOREST_DB_NAME", "gorest_starter_test")
	t.Setenv("GOREST_DB_USER", "integration-user")
	t.Setenv("GOREST_DB_PASSWORD", "integration-password")
	t.Setenv("GOREST_DB_TLS", "verify-full")
	t.Setenv("GOREST_DB_CA_PATH", "")
	t.Setenv("GOREST_DB_CERT_PATH", "/test/client.crt")
	t.Setenv("GOREST_DB_KEY_PATH", "/test/client.key")
	t.Setenv("GOREST_TEST_POSTGRES_CONFIRM_DATABASE", "gorest_starter_test")

	config, err := loadPostgreSQLTestConfig()
	if err != nil {
		t.Fatalf("loadPostgreSQLTestConfig() with server TLS error = %v", err)
	}
	if config.tls != "verify-full" || config.database != "gorest_starter_test" ||
		config.certPath != "/test/client.crt" || config.keyPath != "/test/client.key" {
		t.Fatalf("PostgreSQL test configuration = (TLS:%q, database:%q, cert:%q, key:%q), want verified mTLS and test database",
			config.tls, config.database, config.certPath, config.keyPath)
	}

	environment := config.appEnvironment("3000", t.TempDir(), "test-signature-key")
	if environment["GOREST_DB_DRIVER"] != "postgres" ||
		environment["GOREST_DB_TLS"] != "verify-full" ||
		environment["GOREST_DB_MIGRATION"] != "true" ||
		environment["GOREST_DB_SEEDER"] != "true" {
		t.Fatalf("application database environment does not enable PostgreSQL TLS, migration, and seeder")
	}
	if environment["GOREST_DB_CERT_PATH"] != config.certPath || environment["GOREST_DB_KEY_PATH"] != config.keyPath {
		t.Fatal("PostgreSQL application environment did not preserve the client certificate and key")
	}
}

func TestLoadPostgreSQLTestConfigRequiresExplicitConfirmationForNonTestDatabase(t *testing.T) {
	t.Setenv("GOREST_APP_ENV", "testing")
	t.Setenv("GOREST_DB_DISABLE", "false")
	t.Setenv("GOREST_DB_DRIVER", "postgres")
	t.Setenv("GOREST_DB_HOST", "postgres.example.test")
	t.Setenv("GOREST_DB_NAME", "gorest_starter")
	t.Setenv("GOREST_DB_USER", "integration-user")
	t.Setenv("GOREST_DB_PASSWORD", "integration-password")
	t.Setenv("GOREST_DB_TLS", "verify-full")
	t.Setenv("GOREST_TEST_POSTGRES_CONFIRM_DATABASE", "")
	t.Setenv("GOREST_TEST_POSTGRES_ALLOW_NON_TEST_DATABASE", "")

	if _, err := loadPostgreSQLTestConfig(); err == nil {
		t.Fatal("loadPostgreSQLTestConfig() accepted a non-test database without explicit confirmation")
	}
}

func TestLoadPostgreSQLTestConfigRequiresClientCertificateAndKey(t *testing.T) {
	t.Setenv("GOREST_APP_ENV", "testing")
	t.Setenv("GOREST_DB_DISABLE", "false")
	t.Setenv("GOREST_DB_DRIVER", "postgres")
	t.Setenv("GOREST_DB_HOST", "postgres.example.test")
	t.Setenv("GOREST_DB_NAME", "gorest_starter_test")
	t.Setenv("GOREST_DB_USER", "integration-user")
	t.Setenv("GOREST_DB_PASSWORD", "integration-password")
	t.Setenv("GOREST_DB_TLS", "verify-full")
	t.Setenv("GOREST_DB_CERT_PATH", "/test/client.pem")
	t.Setenv("GOREST_DB_KEY_PATH", "")
	t.Setenv("GOREST_DB_CA_PATH", "")
	t.Setenv("GOREST_TEST_POSTGRES_CONFIRM_DATABASE", "gorest_starter_test")

	if _, err := loadPostgreSQLTestConfig(); err == nil {
		t.Fatal("loadPostgreSQLTestConfig() accepted an incomplete client certificate/key pair")
	}
}
