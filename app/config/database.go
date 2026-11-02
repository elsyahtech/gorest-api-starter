package config

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest/database"
)

// Database returns GOREST's SQLite defaults with this starter's migrations and
// seeders enabled so the example routes are ready after local startup.
func Database() database.Config {
	cfg := database.ConfigDefault
	cfg.DatabaseDirectory = "./app/database"
	// SQLite permits concurrent readers but only one writer at a time. Keep one
	// connection so writes queue and connection-local foreign-key enforcement
	// remains enabled instead of surfacing SQLITE_BUSY or accepting invalid rows.
	// See https://github.com/elsyahtech/gorest-api-starter/blob/main/docs/configuration.md
	cfg.MaxOpenConnections = 1
	cfg.MaxIdleConnections = 1
	cfg.ConnectionMaxLifetime = 0
	cfg.Migration = true
	cfg.Seeder = true

	return cfg
}

func EnableSQLiteForeignKeys(app *gorest.App) error {
	if app == nil || app.Database() == nil {
		return fmt.Errorf("enable SQLite foreign keys: database is not initialized")
	}

	ctx, cancel := app.NewContext(5 * time.Second)
	defer cancel()

	_, message, status, err := app.Database().ExecSQL(ctx, "PRAGMA foreign_keys = ON")
	if err != nil {
		return fmt.Errorf("enable SQLite foreign keys (HTTP %d, %s): %w", status, message, err)
	}

	rows, message, status, err := app.Database().QuerySQL(ctx, "PRAGMA foreign_keys")
	if err != nil {
		return fmt.Errorf("verify SQLite foreign keys (HTTP %d, %s): %w", status, message, err)
	}
	enabled, err := readSQLiteForeignKeySetting(rows)
	if err != nil {
		return err
	}
	if enabled != 1 {
		return fmt.Errorf("verify SQLite foreign keys: pragma value = %d, want 1", enabled)
	}

	return nil
}

func readSQLiteForeignKeySetting(rows *sql.Rows) (enabled int, resultErr error) {
	defer func() {
		if err := rows.Close(); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("close SQLite foreign-key setting rows: %w", err))
		}
	}()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return 0, fmt.Errorf("iterate SQLite foreign-key setting: %w", err)
		}
		return 0, fmt.Errorf("verify SQLite foreign keys: pragma returned no rows")
	}
	if err := rows.Scan(&enabled); err != nil {
		return 0, fmt.Errorf("read SQLite foreign-key setting: %w", err)
	}
	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("iterate SQLite foreign-key setting: %w", err)
	}

	return enabled, nil
}

// LoadDatabase applies environment overrides and rejects unsafe production
// defaults before any application services are registered.
func LoadDatabase() (database.Config, error) {
	if err := validateEnvironment(); err != nil {
		return database.Config{}, err
	}

	cfg := Database()
	production := Environment() == gorest.Production

	if production {
		cfg.Migration = false
		cfg.Seeder = false
	}

	cfg.Driver = strings.ToLower(envString("GOREST_DB_DRIVER", cfg.Driver))

	switch cfg.Driver {
	case database.MYSQL, database.POSTGRES, database.SQLSERVER, database.ORACLE,
		database.SQLITE, database.MONGO, database.SCYLLA:
	default:
		return database.Config{}, fmt.Errorf("GOREST_DB_DRIVER has unsupported value %q", cfg.Driver)
	}

	if cfg.Driver != database.SQLITE {
		cfg.MaxOpenConnections = database.ConfigDefault.MaxOpenConnections
		cfg.MaxIdleConnections = database.ConfigDefault.MaxIdleConnections
	}

	cfg.Host = envString("GOREST_DB_HOST", cfg.Host)
	cfg.Port = envString("GOREST_DB_PORT", cfg.Port)

	if cfg.Port != "" {
		port, err := strconv.Atoi(cfg.Port)

		if err != nil || port < 1 || port > 65535 {
			return database.Config{}, fmt.Errorf("GOREST_DB_PORT must be an integer between 1 and 65535")
		}
	}

	var err error

	cfg.Name = envString("GOREST_DB_NAME", cfg.Name)
	cfg.Username = envString("GOREST_DB_USERNAME", cfg.Username)
	cfg.Password = envSecret("GOREST_DB_PASSWORD", cfg.Password)
	cfg.TLS = envString("GOREST_DB_TLS", cfg.TLS)
	if cfg.AllowInsecureProduction, err = envBool("GOREST_DB_ALLOW_INSECURE_PRODUCTION", cfg.AllowInsecureProduction); err != nil {
		return database.Config{}, err
	}
	cfg.LibDir = envString("GOREST_DB_LIB_DIR", cfg.LibDir)
	cfg.SSLCAPath = envString("GOREST_DB_SSL_CA_PATH", cfg.SSLCAPath)
	cfg.SSLCertPath = envString("GOREST_DB_SSL_CERT_PATH", cfg.SSLCertPath)
	cfg.SSLKeyPath = envString("GOREST_DB_SSL_KEY_PATH", cfg.SSLKeyPath)
	cfg.DatabaseDirectory = envString("GOREST_DB_DIRECTORY", cfg.DatabaseDirectory)

	if cfg.Migration, err = envBool("GOREST_DB_MIGRATION", cfg.Migration); err != nil {
		return database.Config{}, err
	}

	if cfg.Seeder, err = envBool("GOREST_DB_SEEDER", cfg.Seeder); err != nil {
		return database.Config{}, err
	}

	if production {
		if err := validateProductionDatabase(cfg); err != nil {
			return database.Config{}, err
		}
	}

	if cfg.AllowInsecureProduction && production {
		return database.Config{}, fmt.Errorf("GOREST_DB_ALLOW_INSECURE_PRODUCTION must be false in production")
	}

	return cfg, nil
}

func validateProductionDatabase(cfg database.Config) error {
	if err := requireProductionEnv("GOREST_DB_DRIVER"); err != nil {
		return err
	}

	if cfg.Migration || cfg.Seeder {
		return fmt.Errorf("GOREST_DB_MIGRATION and GOREST_DB_SEEDER must be false when APP_ENV=production")
	}

	if cfg.Driver == database.SQLITE {
		if err := requireProductionEnv("GOREST_DB_DIRECTORY"); err != nil {
			return fmt.Errorf("production SQLite requires an explicit database location: %w", err)
		}

		return nil
	}

	keys := []string{
		"GOREST_DB_HOST",
		"GOREST_DB_NAME",
		"GOREST_DB_USERNAME",
		"GOREST_DB_PASSWORD",
	}

	if cfg.Driver != database.SCYLLA {
		keys = append(keys, "GOREST_DB_TLS")
	}

	for _, key := range keys {
		if err := requireProductionEnv(key); err != nil {
			return err
		}
	}

	return nil
}
