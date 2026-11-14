package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/elsyahtech/gorest/database"
	goresthelper "github.com/elsyahtech/gorest/helper"
)

// DatabaseConfig loads the MongoDB connection settings and optional startup tasks.
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
		database.MONGO,
	)
	if !strings.EqualFold(cfg.Driver, database.MONGO) {
		return &database.Config{}, fmt.Errorf("GOREST_DB_DRIVER must be %q for the MongoDB starter", database.MONGO)
	}
	cfg.Driver = database.MONGO

	cfg.Host = strings.TrimSpace(goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_HOST",
		database.DefaultDatabaseHost,
	))
	cfg.Port = strings.TrimSpace(goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_PORT",
		database.MONGO_PORT,
	))
	cfg.Name = strings.TrimSpace(goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_NAME",
		"",
	))
	if cfg.Name == "" {
		return &database.Config{}, errors.New("GOREST_DB_NAME (MongoDB database name) cannot be empty")
	}

	cfg.Username = strings.TrimSpace(goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_USER",
		"",
	))
	cfg.Password = strings.TrimSpace(goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_PASSWORD",
		"",
	))
	if (cfg.Username == "") != (cfg.Password == "") {
		return &database.Config{}, errors.New("GOREST_DB_USER and GOREST_DB_PASSWORD must be configured together")
	}
	cfg.AuthSource = strings.TrimSpace(goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_AUTH_SOURCE",
		"admin",
	))
	if cfg.AuthSource == "" {
		cfg.AuthSource = "admin"
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
		database.TLSFalse,
	)))
	if cfg.TLS != database.TLSFalse && cfg.TLS != "true" {
		return &database.Config{}, errors.New("GOREST_DB_TLS must be true or false for MongoDB")
	}
	cfg.SSLCAPath = strings.TrimSpace(goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_CA_PATH",
		"",
	))
	cfg.SSLCertPath = strings.TrimSpace(goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_CERT_PATH",
		"",
	))
	cfg.SSLKeyPath = strings.TrimSpace(goresthelper.ReadEnvironmentStringWithDefault(
		"GOREST_DB_KEY_PATH",
		"",
	))
	if (cfg.SSLCertPath == "") != (cfg.SSLKeyPath == "") {
		return &database.Config{}, errors.New("GOREST_DB_CERT_PATH and GOREST_DB_KEY_PATH must be configured together")
	}
	if cfg.SSLCertPath != "" && cfg.TLS != "true" {
		return &database.Config{}, errors.New("GOREST_DB_TLS=true is required when MongoDB client certificates are configured")
	}

	cfg.Timeout, err = goresthelper.ReadEnvironmentDurationWithDefault(
		"GOREST_DB_TIMEOUT",
		database.DefaultTimeout,
	)
	if err != nil {
		return &database.Config{}, fmt.Errorf("failed to parse GOREST_DB_TIMEOUT: %w", err)
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
