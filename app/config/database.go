package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/elsyahtech/gorest/database"
)

// DatabaseConfig loads the MySQL connection settings and optional startup tasks.
func DatabaseConfig() (database.Config, error) {
	// Load default configuration
	cfg := database.ConfigDefault

	cfg.Disable = strToBool(os.Getenv("GOREST_DB_DISABLE"), database.DatabaseDisabled)
	cfg.Driver = str(os.Getenv("GOREST_DB_DRIVER"), database.MYSQL)
	cfg.Host = str(os.Getenv("GOREST_DB_HOST"), database.DefaultDatabaseHost)
	cfg.Port = str(os.Getenv("GOREST_DB_PORT"), database.MYSQL_PORT)
	cfg.Name = str(os.Getenv("GOREST_DB_NAME"), database.DefaultDatabaseName)

	cfg.Username = str(os.Getenv("GOREST_DB_USER"), "")
	if cfg.Username == "" {
		return database.Config{}, errors.New("DB Username cannot be empty.")
	}

	cfg.Password = str(os.Getenv("GOREST_DB_PASSWORD"), "")
	if cfg.Password == "" {
		return database.Config{}, fmt.Errorf("DB Password cannot be empty.")
	}

	cfg.AllowInsecureProduction = strToBool(os.Getenv("GOREST_DB_ALLOW_INSECURE_PRODUCTION"), false)
	cfg.TLS = str(os.Getenv("GOREST_DB_TLS"), database.TLSDisabled)

	cfg.SSLCertPath = str(os.Getenv("GOREST_DB_CERT_PATH"), "")
	cfg.SSLCAPath = str(os.Getenv("GOREST_DB_CA_PATH"), "")
	if cfg.Driver == database.MYSQL && cfg.TLS == database.TLSTrue && cfg.SSLCAPath == "" {
		return database.Config{}, fmt.Errorf("GOREST_DB_CA_PATH cannot be empty when MySQL TLS is enabled")
	}

	cfg.SSLKeyPath = str(os.Getenv("GOREST_DB_KEY_PATH"), "")
	cfg.Charset = str(os.Getenv("GOREST_DB_CHARSET"), database.DefaultCharset)
	cfg.Timeout = strToTimeDuration(os.Getenv("GOREST_DB_TIMEOUT"), database.DefaultTimeout)
	cfg.ConnectionMaxLifetime = strToTimeDuration(os.Getenv("GOREST_DB_CONNECTION_MAX_LIFETIME"), database.DefaultConnectionMaxLifetime)
	cfg.MaxOpenConnections = strToInt(os.Getenv("GOREST_DB_MAX_OPEN_CONNECTIONS"), database.DefaultMaxOpenConnections)
	cfg.MaxIdleConnections = strToInt(os.Getenv("GOREST_DB_MAX_IDLE_CONNECTIONS"), database.DefaultMaxIdleConnections)
	cfg.QueryLogging = strToBool(os.Getenv("GOREST_DB_QUERY_LOGGING"), database.DefaultQueryLogging)
	cfg.DatabaseDirectory = str(os.Getenv("GOREST_DB_DIRECTORY"), database.DefaultDatabaseDirectory)
	cfg.Migration = strToBool(os.Getenv("GOREST_DB_MIGRATION"), database.MigrationDisabled)
	cfg.Seeder = strToBool(os.Getenv("GOREST_DB_SEEDER"), database.SeederDisabled)
	cfg.MigrationDirectory = str(os.Getenv("GOREST_DB_MIGRATIONS_DIRECTORY"), database.DefaultMigrationDirectory)
	cfg.SeederDirectory = str(os.Getenv("GOREST_DB_SEEDERS_DIRECTORY"), database.DefaultSeederDirectory)

	return cfg, nil
}
