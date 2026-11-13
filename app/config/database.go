package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/elsyahtech/gorest/database"
	goresthelper "github.com/elsyahtech/gorest/helper"
)

// DatabaseConfig loads the PostgreSQL connection settings and optional startup tasks.
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

	// This starter supports PostgreSQL only.
	cfg.Driver = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_DRIVER",
		database.POSTGRES,
	)
	if !strings.EqualFold(cfg.Driver, database.POSTGRES) {
		return &database.Config{}, fmt.Errorf("GOREST_DB_DRIVER must be %q for the PostgreSQL starter", database.POSTGRES)
	}

	cfg.Driver = database.POSTGRES

	cfg.Host = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_HOST",
		database.DefaultDatabaseHost,
	)

	cfg.Port = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_PORT",
		database.POSTGRES_PORT,
	)

	cfg.Name = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_NAME",
		database.DefaultDatabaseName,
	)

	cfg.Username = strings.TrimSpace(goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_USER",
		"",
	))
	if cfg.Username == "" {
		return &database.Config{}, errors.New("DB Username cannot be empty")
	}

	cfg.Password = strings.TrimSpace(goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_PASSWORD",
		"",
	))
	if cfg.Password == "" {
		return &database.Config{}, errors.New("DB Password cannot be empty")
	}

	cfg.AllowInsecureProduction, err = goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_DB_ALLOW_INSECURE_PRODUCTION",
		false,
	)
	if err != nil {
		return &database.Config{}, fmt.Errorf("failed to parse GOREST_DB_ALLOW_INSECURE_PRODUCTION: %w", err)
	}

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

	cfg.TLS = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_TLS",
		database.TLSDisabled,
	)

	cfg.SSLCertPath = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_CERT_PATH",
		"",
	)

	cfg.SSLCAPath = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_CA_PATH",
		"",
	)

	cfg.SSLKeyPath = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_KEY_PATH",
		"",
	)

	cfg.Charset = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_CHARSET",
		database.DefaultCharset,
	)

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
		database.DefaultDatabaseDirectory,
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
