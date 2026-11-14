//go:build integration

package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/elsyahtech/gorest-api-starter/app/config"
	"github.com/elsyahtech/gorest/database"
)

func TestScyllaDatabaseConnection(t *testing.T) {
	testConfig, err := loadScyllaTestConfig()
	if err != nil {
		t.Fatal(err)
	}
	applyScyllaTestConfig(t, testConfig)

	dbConfig, err := config.DatabaseConfig()
	if err != nil {
		t.Fatalf("load ScyllaDB connection-test configuration: %v", err)
	}
	if dbConfig.Disable || dbConfig.Driver != database.SCYLLA ||
		dbConfig.Name != testConfig.keyspace || dbConfig.TLS != testConfig.tls {
		t.Fatal("database configuration does not match the confirmed ScyllaDB test keyspace")
	}
	dbConfig.Migration = false
	dbConfig.Seeder = false
	dbConfig.Timeout = 20 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), dbConfig.Timeout)
	defer cancel()
	db, message, err := database.New(ctx, dbConfig, "UTC")
	if err != nil {
		t.Fatalf("connect to configured ScyllaDB test keyspace (%s): %v", message, err)
	}
	t.Cleanup(func() {
		closeScyllaDatabase(t, db, dbConfig)
	})
	if message, err := db.Ping(ctx, dbConfig); err != nil {
		t.Fatalf("verify ScyllaDB test keyspace connection (%s): %v", message, err)
	}

	var serverVersion string
	if err := db.Scylla.Query("SELECT release_version FROM system.local").
		WithContext(ctx).Scan(&serverVersion); err != nil {
		t.Fatalf("query ScyllaDB server version: %v", err)
	}
	if serverVersion == "" {
		t.Fatal("ScyllaDB returned an empty server version")
	}
	t.Logf("Connected to ScyllaDB server version %s", serverVersion)
}
