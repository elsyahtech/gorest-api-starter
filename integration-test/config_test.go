//go:build integration

package integration_test

import "testing"

func TestLoadMySQLTestConfigAllowsTLSWithoutCAPath(t *testing.T) {
	t.Setenv("GOREST_TEST_MYSQL_HOST", "mysql.example.test")
	t.Setenv("GOREST_TEST_MYSQL_DATABASE", "gorest_starter_test")
	t.Setenv("GOREST_TEST_MYSQL_USER", "integration-user")
	t.Setenv("GOREST_TEST_MYSQL_PASSWORD", "integration-password")
	t.Setenv("GOREST_TEST_MYSQL_TLS", "true")
	t.Setenv("GOREST_TEST_MYSQL_CONFIRM_DATABASE", "gorest_starter_test")

	config, err := loadMySQLTestConfig()
	if err != nil {
		t.Fatalf("loadMySQLTestConfig() with TLS and no CA path error = %v", err)
	}
	if config.tls != "true" {
		t.Errorf("TLS mode = %q, want true", config.tls)
	}

	environment := config.appEnvironment("3000", t.TempDir(), "test-signature-key")
	if environment["GOREST_DB_TLS"] != "true" {
		t.Errorf("application TLS mode = %q, want true", environment["GOREST_DB_TLS"])
	}
	if environment["GOREST_DB_CA_PATH"] != "" {
		t.Errorf("application CA path = %q, want empty", environment["GOREST_DB_CA_PATH"])
	}
}
