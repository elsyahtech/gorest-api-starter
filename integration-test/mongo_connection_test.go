//go:build integration

package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/elsyahtech/gorest-api-starter/app/config"
	"github.com/elsyahtech/gorest/database"
)

func TestMongoDatabaseConnection(t *testing.T) {
	testConfig, err := loadMongoTestConfig()
	if err != nil {
		t.Fatal(err)
	}
	applyMongoTestConfig(t, testConfig)

	dbConfig, err := config.DatabaseConfig()
	if err != nil {
		t.Fatalf("load MongoDB connection-test configuration: %v", err)
	}
	if dbConfig.Disable || dbConfig.Driver != database.MONGO ||
		dbConfig.Name != testConfig.database || dbConfig.AuthSource != testConfig.authSource ||
		dbConfig.TLS != testConfig.tls {
		t.Fatal("database configuration does not match the confirmed MongoDB test database")
	}
	dbConfig.Migration = false
	dbConfig.Seeder = false
	dbConfig.Timeout = 20 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), dbConfig.Timeout)
	defer cancel()
	db, message, err := database.New(ctx, dbConfig, "UTC")
	if err != nil {
		t.Fatalf("connect to configured MongoDB test database (%s): %v", message, err)
	}
	t.Cleanup(func() {
		cleanupMongoTestConnection(t, db.Mongo)
	})
	if err := db.Mongo.Ping(ctx, nil); err != nil {
		t.Fatalf("verify MongoDB test database connection: %v", err)
	}
}
