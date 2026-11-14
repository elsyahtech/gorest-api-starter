//go:build integration

package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/elsyahtech/gorest-api-starter/app/config"
	"github.com/elsyahtech/gorest/database"
)

func TestSQLServerTLSConnection(t *testing.T) {
	testConfig, err := loadSQLServerTestConfig()
	if err != nil {
		t.Fatal(err)
	}
	applySQLServerTestConfig(t, testConfig)

	dbConfig, err := config.DatabaseConfig()
	if err != nil {
		t.Fatalf("load SQL Server TLS configuration: %v", err)
	}
	if dbConfig.Disable || dbConfig.Driver != database.SQLSERVER ||
		dbConfig.TLS != testConfig.tls || dbConfig.TrustServerCertificate {
		t.Fatal("database configuration does not require verified SQL Server TLS")
	}
	dbConfig.Migration = false
	dbConfig.Seeder = false
	dbConfig.Timeout = 15 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), dbConfig.Timeout)
	defer cancel()
	db, message, err := database.New(ctx, dbConfig, "UTC")
	if err != nil {
		t.Fatalf("connect to configured SQL Server with verified TLS (%s): %v", message, err)
	}
	t.Cleanup(func() {
		if err := db.SQL.Close(); err != nil {
			t.Errorf("close SQL Server TLS integration connection: %v", err)
		}
	})
}
