package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/elsyahtech/gorest/database"
)

// DatabaseConfig loads the PostgreSQL connection settings and optional startup tasks.
func DatabaseConfig() (*database.Config, error) {
	cfg := database.ConfigDefault

	var err error

	cfg.Disable, err = envStrToBool("GOREST_DB_DISABLE", database.DatabaseDisabled)
	if err != nil {
		return &database.Config{}, fmt.Errorf("failed to parse GOREST_DB_DISABLE: %w", err)
	}

	// This starter supports PostgreSQL only.
	cfg.Driver = envStr("GOREST_DB_DRIVER", database.POSTGRES)
	if !strings.EqualFold(cfg.Driver, database.POSTGRES) {
		return &database.Config{}, fmt.Errorf("GOREST_DB_DRIVER must be %q for the PostgreSQL starter", database.POSTGRES)
	}

	cfg.Driver = database.POSTGRES
	cfg.Host = envStr("GOREST_DB_HOST", database.DefaultDatabaseHost)
	cfg.Port = envStr("GOREST_DB_PORT", database.POSTGRES_PORT)
	cfg.Name = envStr("GOREST_DB_NAME", database.DefaultDatabaseName)

	cfg.Username = envStr("GOREST_DB_USER", "")
	if cfg.Username == "" {
		return &database.Config{}, errors.New("DB Username cannot be empty")
	}

	cfg.Password = envStr("GOREST_DB_PASSWORD", "")
	if cfg.Password == "" {
		return &database.Config{}, errors.New("DB Password cannot be empty")
	}

	cfg.AllowInsecureProduction, err = envStrToBool("GOREST_DB_ALLOW_INSECURE_PRODUCTION", false)
	if err != nil {
		return &database.Config{}, fmt.Errorf("failed to parse GOREST_DB_ALLOW_INSECURE_PRODUCTION: %w", err)
	}

	timeout, err := envStrToTimeDuration("GOREST_DB_TIMEOUT", database.DefaultTimeout)
	if err != nil {
		return &database.Config{}, fmt.Errorf("%w", err)
	}

	connectionMaxLifetime, err := envStrToTimeDuration("GOREST_DB_CONNECTION_MAX_LIFETIME", database.DefaultConnectionMaxLifetime)
	if err != nil {
		return &database.Config{}, fmt.Errorf("%w", err)
	}

	cfg.TLS = envStr("GOREST_DB_TLS", database.TLSDisabled)
	cfg.SSLCertPath = envStr("GOREST_DB_CERT_PATH", "")
	cfg.SSLCAPath = envStr("GOREST_DB_CA_PATH", "")
	cfg.SSLKeyPath = envStr("GOREST_DB_KEY_PATH", "")
	cfg.Charset = envStr("GOREST_DB_CHARSET", database.DefaultCharset)
	cfg.Timeout = timeout
	cfg.ConnectionMaxLifetime = connectionMaxLifetime
	cfg.MaxOpenConnections = envStrToInt("GOREST_DB_MAX_OPEN_CONNECTIONS", database.DefaultMaxOpenConnections)
	cfg.MaxIdleConnections = envStrToInt("GOREST_DB_MAX_IDLE_CONNECTIONS", database.DefaultMaxIdleConnections)

	cfg.QueryLogging, err = envStrToBool("GOREST_DB_QUERY_LOGGING", database.DefaultQueryLogging)
	if err != nil {
		return &database.Config{}, fmt.Errorf("failed to parse GOREST_DB_QUERY_LOGGING: %w", err)
	}

	cfg.DatabaseDirectory = envStr("GOREST_DB_DIRECTORY", database.DefaultDatabaseDirectory)

	cfg.Migration, err = envStrToBool("GOREST_DB_MIGRATION", database.MigrationDisabled)
	if err != nil {
		return &database.Config{}, fmt.Errorf("failed to parse GOREST_DB_MIGRATION: %w", err)
	}

	cfg.Seeder, err = envStrToBool("GOREST_DB_SEEDER", database.SeederDisabled)
	if err != nil {
		return &database.Config{}, fmt.Errorf("failed to parse GOREST_DB_SEEDER: %w", err)
	}

	cfg.MigrationDirectory = envStr("GOREST_DB_MIGRATIONS_DIRECTORY", database.DefaultMigrationDirectory)
	cfg.SeederDirectory = envStr("GOREST_DB_SEEDERS_DIRECTORY", database.DefaultSeederDirectory)

	return &cfg, nil
}
