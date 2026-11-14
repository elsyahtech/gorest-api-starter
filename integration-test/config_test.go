//go:build integration

package integration_test

import "testing"

func configureSQLServerTestEnvironment(t *testing.T) {
	t.Helper()
	t.Setenv("GOREST_APP_ENV", "testing")
	t.Setenv("GOREST_DB_DISABLE", "false")
	t.Setenv("GOREST_DB_DRIVER", "sqlserver")
	t.Setenv("GOREST_DB_HOST", "sqlserver.example.test")
	t.Setenv("GOREST_DB_PORT", "")
	t.Setenv("GOREST_DB_NAME", "gorest-dedicated-db")
	t.Setenv("GOREST_DB_USER", "integration-user")
	t.Setenv("GOREST_DB_PASSWORD", "integration-password")
	t.Setenv("GOREST_DB_TLS", "mandatory")
	t.Setenv("GOREST_DB_CA_PATH", "")
	t.Setenv("GOREST_DB_TRUST_SERVER_CERTIFICATE", "false")
	t.Setenv("GOREST_TEST_SQLSERVER_CONFIRM_DATABASE", "gorest-dedicated-db")
}

func TestLoadSQLServerTestConfigAcceptsVerifiedTLS(t *testing.T) {
	configureSQLServerTestEnvironment(t)

	config, err := loadSQLServerTestConfig()
	if err != nil {
		t.Fatalf("loadSQLServerTestConfig() with verified TLS error = %v", err)
	}
	if config.tls != "mandatory" || config.database != "gorest-dedicated-db" ||
		config.port != "1433" {
		t.Fatalf("SQL Server test configuration = (TLS:%q, database-set:%t, port:%q), want mandatory TLS, confirmed database, and default port",
			config.tls, config.database != "", config.port)
	}

	environment := config.appEnvironment("3000", t.TempDir(), "test-signature-key")
	if environment["GOREST_DB_DRIVER"] != "sqlserver" ||
		environment["GOREST_DB_TLS"] != "mandatory" ||
		environment["GOREST_DB_CA_PATH"] != "" ||
		environment["GOREST_DB_TRUST_SERVER_CERTIFICATE"] != "false" ||
		environment["GOREST_DB_MIGRATION"] != "true" ||
		environment["GOREST_DB_SEEDER"] != "true" {
		t.Fatal("application environment does not enforce verified SQL Server TLS and enable migration/seeder")
	}
}

func TestLoadSQLServerTestConfigRequiresExactDatabaseConfirmation(t *testing.T) {
	configureSQLServerTestEnvironment(t)
	t.Setenv("GOREST_TEST_SQLSERVER_CONFIRM_DATABASE", "another-database")

	if _, err := loadSQLServerTestConfig(); err == nil {
		t.Fatal("loadSQLServerTestConfig() accepted a database without exact-name confirmation")
	}
}

func TestLoadSQLServerTestConfigRejectsSystemDatabase(t *testing.T) {
	configureSQLServerTestEnvironment(t)
	t.Setenv("GOREST_DB_NAME", "master")
	t.Setenv("GOREST_TEST_SQLSERVER_CONFIRM_DATABASE", "master")

	if _, err := loadSQLServerTestConfig(); err == nil {
		t.Fatal("loadSQLServerTestConfig() accepted a SQL Server system database")
	}
}

func TestLoadSQLServerTestConfigRejectsUnverifiedServerCertificate(t *testing.T) {
	configureSQLServerTestEnvironment(t)
	t.Setenv("GOREST_DB_TRUST_SERVER_CERTIFICATE", "true")

	if _, err := loadSQLServerTestConfig(); err == nil {
		t.Fatal("loadSQLServerTestConfig() accepted TrustServerCertificate=true")
	}
}

func TestLoadSQLServerTestConfigRequiresEncryptedTLS(t *testing.T) {
	configureSQLServerTestEnvironment(t)
	t.Setenv("GOREST_DB_TLS", "disable")

	if _, err := loadSQLServerTestConfig(); err == nil {
		t.Fatal("loadSQLServerTestConfig() accepted disabled SQL Server TLS")
	}
}
