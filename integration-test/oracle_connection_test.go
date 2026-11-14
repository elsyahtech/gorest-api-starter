//go:build integration

package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/elsyahtech/gorest-api-starter/app/config"
	"github.com/elsyahtech/gorest/database"
)

func TestOracleDatabaseConnection(t *testing.T) {
	testConfig, err := loadOracleTestConfig()
	if err != nil {
		t.Fatal(err)
	}
	applyOracleTestConfig(t, testConfig)

	dbConfig, err := config.DatabaseConfig()
	if err != nil {
		t.Fatalf("load Oracle connection-test configuration: %v", err)
	}
	if dbConfig.Disable || dbConfig.Driver != database.ORACLE ||
		dbConfig.Name != testConfig.serviceName || dbConfig.Username != testConfig.username ||
		dbConfig.TLS != testConfig.tls || dbConfig.LibDir != testConfig.libDir ||
		dbConfig.SSLCAPath != testConfig.walletPath {
		t.Fatal("database configuration does not match the confirmed Oracle test service and schema")
	}
	dbConfig.Migration = false
	dbConfig.Seeder = false
	dbConfig.Timeout = 20 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), dbConfig.Timeout)
	defer cancel()
	db, message, err := database.New(ctx, dbConfig, "UTC")
	if err != nil {
		t.Fatalf("connect to configured Oracle test service (%s): %v", message, err)
	}
	t.Cleanup(func() {
		if err := db.SQL.Close(); err != nil {
			t.Errorf("close Oracle test connection: %v", err)
		}
	})

	var result int
	if err := db.SQL.QueryRowContext(ctx, "SELECT 1 FROM dual").Scan(&result); err != nil {
		t.Fatalf("verify Oracle test connection: %v", err)
	}
	if result != 1 {
		t.Fatalf("Oracle test connection query returned %d, want 1", result)
	}
}
