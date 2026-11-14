package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/elsyahtech/gorest/database"
	goresthelper "github.com/elsyahtech/gorest/helper"
)

// DatabaseConfig loads the Oracle connection settings and optional startup tasks.
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

	cfg.Driver = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_DRIVER",
		database.ORACLE,
	)
	if !strings.EqualFold(cfg.Driver, database.ORACLE) {
		return &database.Config{}, fmt.Errorf("GOREST_DB_DRIVER must be %q for the Oracle starter", database.ORACLE)
	}
	cfg.Driver = database.ORACLE

	cfg.Host = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_HOST",
		database.DefaultDatabaseHost,
	)
	cfg.Port = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_PORT",
		database.ORACLE_PORT,
	)
	cfg.LibDir = goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_LIB_DIR_PATH",
		database.DefaultDatabaseLibDir,
	)
	cfg.Name = strings.TrimSpace(goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_NAME",
		"",
	))
	if cfg.Name == "" {
		return &database.Config{}, errors.New("GOREST_DB_NAME (Oracle service name) cannot be empty")
	}

	cfg.Username = strings.TrimSpace(goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_USER",
		"",
	))
	if cfg.Username == "" {
		return &database.Config{}, errors.New("GOREST_DB_USER cannot be empty")
	}
	cfg.Password = strings.TrimSpace(goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_PASSWORD",
		"",
	))
	if cfg.Password == "" {
		return &database.Config{}, errors.New("GOREST_DB_PASSWORD cannot be empty")
	}

	cfg.AllowInsecureProduction, err = goresthelper.ReadEnvironmentBooleanWithDefault(
		"GOREST_DB_ALLOW_INSECURE_PRODUCTION",
		false,
	)
	if err != nil {
		return &database.Config{}, fmt.Errorf("failed to parse GOREST_DB_ALLOW_INSECURE_PRODUCTION: %w", err)
	}

	cfg.TLS = strings.ToLower(strings.TrimSpace(goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_TLS",
		database.TLSDisabled,
	)))
	if cfg.TLS != database.TLSEnabled && cfg.TLS != database.TLSDisabled {
		return &database.Config{}, fmt.Errorf("GOREST_DB_TLS must be %q or %q for Oracle", database.TLSEnabled, database.TLSDisabled)
	}
	cfg.SSLCAPath = strings.TrimSpace(goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_CA_PATH",
		database.DefaultSSLCAPath,
	))
	if cfg.TLS == database.TLSEnabled && cfg.SSLCAPath == "" {
		return &database.Config{}, errors.New("GOREST_DB_CA_PATH must point to an Oracle wallet directory when GOREST_DB_TLS=enable")
	}
	cfg.SSLCertPath = ""
	cfg.SSLKeyPath = ""
	cfg.TrustServerCertificate = false

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
