//go:build integration

package integration_test

import (
	"path/filepath"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/config"
	"github.com/elsyahtech/gorest/database"
)

func TestSQLiteTestConfigUsesTemporaryDatabaseAndRepositoryMigrations(t *testing.T) {
	testConfig := newSQLiteTestConfig(t)
	applySQLiteTestConfig(t, testConfig)

	dbConfig, err := config.DatabaseConfig()
	if err != nil {
		t.Fatalf("load SQLite test configuration: %v", err)
	}
	resolved := database.Run(*dbConfig)
	projectRoot, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("resolve repository root: %v", err)
	}

	if resolved.Driver != database.SQLITE ||
		resolved.DatabaseDirectory != testConfig.databaseDirectory ||
		resolved.MigrationDirectory != filepath.Join(projectRoot, "app", "database", "migrations", "sqlite") ||
		resolved.SeederDirectory != filepath.Join(projectRoot, "app", "database", "seeders", "sqlite") {
		t.Fatalf("resolved SQLite test config has incorrect driver or asset directories")
	}
	if !filepath.IsAbs(resolved.DatabaseDirectory) {
		t.Fatal("SQLite integration database directory is not isolated in an absolute temporary directory")
	}
}
