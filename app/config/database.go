package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/elsyahtech/gorest/database"
)

// DatabaseConfig loads the MySQL connection settings and optional startup tasks.
func DatabaseConfig() (database.Config, error) {
	cfg := database.ConfigDefault

	var err error

	cfg.Disable, err = strToBool("GOREST_DB_DISABLE", database.DatabaseDisabled)
	if err != nil {
		return database.Config{}, fmt.Errorf("failed to parse GOREST_DB_DISABLE: %w", err)
	}

	// This starter supports MySQL only.
	cfg.Driver = str("GOREST_DB_DRIVER", database.MYSQL)
	if !strings.EqualFold(cfg.Driver, database.MYSQL) {
		return database.Config{}, fmt.Errorf("GOREST_DB_DRIVER must be %q for the MySQL starter", database.MYSQL)
	}

	cfg.Driver = database.MYSQL
	cfg.Host = str("GOREST_DB_HOST", database.DefaultDatabaseHost)
	cfg.Port = str("GOREST_DB_PORT", database.MYSQL_PORT)
	cfg.Name = str("GOREST_DB_NAME", database.DefaultDatabaseName)

	cfg.Username = str("GOREST_DB_USER", "")
	if cfg.Username == "" {
		return database.Config{}, errors.New("DB Username cannot be empty")
	}

	cfg.Password = str("GOREST_DB_PASSWORD", "")
	if cfg.Password == "" {
		return database.Config{}, errors.New("DB Password cannot be empty")
	}

	cfg.AllowInsecureProduction, err = strToBool("GOREST_DB_ALLOW_INSECURE_PRODUCTION", false)
	if err != nil {
		return database.Config{}, fmt.Errorf("failed to parse GOREST_DB_ALLOW_INSECURE_PRODUCTION: %w", err)
	}

	cfg.TLS = str("GOREST_DB_TLS", database.TLSDisabled)
	cfg.SSLCertPath = str("GOREST_DB_CERT_PATH", "")
	cfg.SSLCAPath = str("GOREST_DB_CA_PATH", "")
	cfg.SSLKeyPath = str("GOREST_DB_KEY_PATH", "")
	cfg.Charset = str("GOREST_DB_CHARSET", database.DefaultCharset)
	cfg.Timeout = strToTimeDuration("GOREST_DB_TIMEOUT", database.DefaultTimeout)
	cfg.ConnectionMaxLifetime = strToTimeDuration("GOREST_DB_CONNECTION_MAX_LIFETIME", database.DefaultConnectionMaxLifetime)
	cfg.MaxOpenConnections = strToInt("GOREST_DB_MAX_OPEN_CONNECTIONS", database.DefaultMaxOpenConnections)
	cfg.MaxIdleConnections = strToInt("GOREST_DB_MAX_IDLE_CONNECTIONS", database.DefaultMaxIdleConnections)

	cfg.QueryLogging, err = strToBool("GOREST_DB_QUERY_LOGGING", database.DefaultQueryLogging)
	if err != nil {
		return database.Config{}, fmt.Errorf("failed to parse GOREST_DB_QUERY_LOGGING: %w", err)
	}

	cfg.DatabaseDirectory = str("GOREST_DB_DIRECTORY", database.DefaultDatabaseDirectory)

	cfg.Migration, err = strToBool("GOREST_DB_MIGRATION", database.MigrationDisabled)
	if err != nil {
		return database.Config{}, fmt.Errorf("failed to parse GOREST_DB_MIGRATION: %w", err)
	}

	cfg.Seeder, err = strToBool("GOREST_DB_SEEDER", database.SeederDisabled)
	if err != nil {
		return database.Config{}, fmt.Errorf("failed to parse GOREST_DB_SEEDER: %w", err)
	}

	cfg.MigrationDirectory = str("GOREST_DB_MIGRATIONS_DIRECTORY", database.DefaultMigrationDirectory)
	cfg.SeederDirectory = str("GOREST_DB_SEEDERS_DIRECTORY", database.DefaultSeederDirectory)

	return cfg, nil
}
