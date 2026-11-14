//go:build integration

package integration_test

import (
	"testing"

	"github.com/elsyahtech/gorest/database"
)

func configureOracleTestEnvironment(t *testing.T) {
	t.Helper()
	t.Setenv("GOREST_APP_ENV", "testing")
	t.Setenv("GOREST_DB_DISABLE", "false")
	t.Setenv("GOREST_DB_DRIVER", database.ORACLE)
	t.Setenv("GOREST_DB_HOST", "oracle.example.test")
	t.Setenv("GOREST_DB_PORT", "")
	t.Setenv("GOREST_DB_NAME", "FREEPDB1")
	t.Setenv("GOREST_DB_USER", "gorest_test")
	t.Setenv("GOREST_DB_PASSWORD", "integration-password")
	t.Setenv("GOREST_DB_TLS", database.TLSDisabled)
	t.Setenv("GOREST_DB_CA_PATH", "")
	t.Setenv("GOREST_DB_LIB_DIR_PATH", "")
	t.Setenv("GOREST_TEST_ORACLE_CONFIRM_USER", "gorest_test")
}

func TestLoadOracleTestConfigAcceptsDedicatedSchema(t *testing.T) {
	configureOracleTestEnvironment(t)

	config, err := loadOracleTestConfig()
	if err != nil {
		t.Fatalf("loadOracleTestConfig() error = %v", err)
	}
	if config.tls != database.TLSDisabled || config.serviceName != "FREEPDB1" ||
		config.username != "gorest_test" || config.port != database.ORACLE_PORT {
		t.Fatalf("Oracle test configuration = (TLS:%q, service:%q, schema:%q, port:%q), want configured service, confirmed schema, and default port",
			config.tls, config.serviceName, config.username, config.port)
	}
	environment := config.appEnvironment("3000", t.TempDir(), "test-signature-key")
	if environment["GOREST_DB_DRIVER"] != database.ORACLE ||
		environment["GOREST_DB_NAME"] != "FREEPDB1" ||
		environment["GOREST_DB_USER"] != "gorest_test" ||
		environment["GOREST_DB_TLS"] != database.TLSDisabled ||
		environment["GOREST_DB_MIGRATION"] != "true" ||
		environment["GOREST_DB_SEEDER"] != "true" {
		t.Fatal("application environment does not select Oracle and enable migrations/seeders")
	}
}

func TestLoadOracleTestConfigRequiresExactSchemaConfirmation(t *testing.T) {
	configureOracleTestEnvironment(t)
	t.Setenv("GOREST_TEST_ORACLE_CONFIRM_USER", "another-user")

	if _, err := loadOracleTestConfig(); err == nil {
		t.Fatal("loadOracleTestConfig() accepted a schema without exact username confirmation")
	}
}

func TestLoadOracleTestConfigRejectsAdministrativeAccounts(t *testing.T) {
	configureOracleTestEnvironment(t)
	t.Setenv("GOREST_DB_USER", "SYS")
	t.Setenv("GOREST_TEST_ORACLE_CONFIRM_USER", "SYS")

	if _, err := loadOracleTestConfig(); err == nil {
		t.Fatal("loadOracleTestConfig() accepted the Oracle SYS account")
	}
}

func TestLoadOracleTestConfigRejectsProduction(t *testing.T) {
	configureOracleTestEnvironment(t)
	t.Setenv("GOREST_APP_ENV", "production")

	if _, err := loadOracleTestConfig(); err == nil {
		t.Fatal("loadOracleTestConfig() accepted production mode")
	}
}

func TestLoadOracleTestConfigRequiresWalletForTCPS(t *testing.T) {
	configureOracleTestEnvironment(t)
	t.Setenv("GOREST_DB_TLS", database.TLSEnabled)

	if _, err := loadOracleTestConfig(); err == nil {
		t.Fatal("loadOracleTestConfig() accepted TCPS without a wallet directory")
	}
}
