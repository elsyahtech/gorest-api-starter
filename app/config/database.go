package config

import (
	"fmt"
	"strings"

	"github.com/elsyahtech/gorest/database"
	goresthelper "github.com/elsyahtech/gorest/helper"
)

// DatabaseConfig loads the SQLite database settings and optional startup tasks.
func DatabaseConfig() (*database.Config, error) {
	cfg := database.ConfigDefault

	var err error

	cfg.Disable, err = goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_DB_DISABLE",
		database.DatabaseDisabled,
	)
	if err != nil {
		return &database.Config{}, fmt.Errorf("failed to parse GOREST_DB_DISABLE: %w", err)
	}

	// This starter supports SQLite only.
	cfg.Driver = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_DRIVER",
		database.SQLITE,
	)
	if !strings.EqualFold(cfg.Driver, database.SQLITE) {
		return &database.Config{}, fmt.Errorf("GOREST_DB_DRIVER must be %q for the SQLite starter", database.SQLITE)
	}

	cfg.Driver = database.SQLITE

	cfg.Name = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_NAME",
		database.DefaultDatabaseName,
	)

	timeout, err := goresthelper.ReadEnvironmentDurationWithDefault(
		"GOREST_DB_TIMEOUT",
		database.DefaultTimeout,
	)
	if err != nil {
		return &database.Config{}, fmt.Errorf("%w", err)
	}

	connectionMaxLifetime, err := goresthelper.ReadEnvironmentDurationWithDefault(
		"GOREST_DB_CONNECTION_MAX_LIFETIME",
		database.DefaultConnectionMaxLifetime,
	)
	if err != nil {
		return &database.Config{}, fmt.Errorf("%w", err)
	}

	cfg.Timeout = timeout
	cfg.ConnectionMaxLifetime = connectionMaxLifetime

	cfg.MaxOpenConnections, err = goresthelper.ReadEnvironmentIntegerWithDefault(
		"GOREST_DB_MAX_OPEN_CONNECTIONS",
		database.DefaultMaxOpenConnections,
	)
	if err != nil {
		return &database.Config{}, fmt.Errorf("failed to parse GOREST_DB_MAX_OPEN_CONNECTIONS: %w", err)
	}

	cfg.MaxIdleConnections, err = goresthelper.ReadEnvironmentIntegerWithDefault(
		"GOREST_DB_MAX_IDLE_CONNECTIONS",
		database.DefaultMaxIdleConnections,
	)
	if err != nil {
		return &database.Config{}, fmt.Errorf("failed to parse GOREST_DB_MAX_IDLE_CONNECTIONS: %w", err)
	}

	cfg.QueryLogging, err = goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_DB_QUERY_LOGGING",
		database.DefaultQueryLogging,
	)
	if err != nil {
		return &database.Config{}, fmt.Errorf("failed to parse GOREST_DB_QUERY_LOGGING: %w", err)
	}

	cfg.DatabaseDirectory = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_DIRECTORY",
		"./app/database",
	)

	cfg.Migration, err = goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_DB_MIGRATION",
		database.MigrationDisabled,
	)
	if err != nil {
		return &database.Config{}, fmt.Errorf("failed to parse GOREST_DB_MIGRATION: %w", err)
	}

	cfg.Seeder, err = goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_DB_SEEDER",
		database.SeederDisabled,
	)
	if err != nil {
		return &database.Config{}, fmt.Errorf("failed to parse GOREST_DB_SEEDER: %w", err)
	}

	cfg.MigrationDirectory = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_MIGRATIONS_DIRECTORY",
		database.DefaultMigrationDirectory,
	)

	cfg.SeederDirectory = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_SEEDERS_DIRECTORY",
		database.DefaultSeederDirectory,
	)

	return &cfg, nil
}
