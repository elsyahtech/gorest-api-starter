package config

import (
	"time"

	"github.com/elsyahtech/gorest/database"
)

// Database is the configuration for the Database connection.
// GOREST connects to Database based on the values you set here, but only when Enabled is true,
// so adjust them to match your environment.
func Database() database.Config {
	return database.Config{
		// Enabled determines whether the Database connection is active and running.
		// Set to true to enable Database integration, or false to disable it.
		// Example: true
		// Default: false
		Enabled: true,

		// Database engine/driver type.
		// Supported: MYSQL, POSTGRES, SQLSERVER, ORACLE, SQLITE, MONGO, SCYLLA
		// Default "sqlite"; if driver not configure
		Driver: database.POSTGRES,

		// Database server port number.
		// Common defaults: MySQL: 3306, PostgreSQL: 5432, MSSQL: 1433, MongoDB: 27017, ScyllaDB: 9042, Oracle: 1521
		// Default: empty
		Port: database.POSTGRES_PORT,

		// Database server hostname or IP address.
		// Development: use "localhost" or "127.0.0.1"
		// Production: use DNS name or private IP
		// Default 127.0.0.1
		Host: "127.0.0.1",

		// Database login username.
		// Development defaults: MySQL: root, PostgreSQL: postgres, MSSQL: sa, MongoDB: admin, Oracle: system
		// Production: Create a dedicated app user with minimal permissions.
		// No default
		Username: "",

		// Database login password. This value is compiled into the application binary.
		// Anyone with access to the binary may be able to recover it; protect the binary
		// and use a deployment-specific credential. Do not commit a real production secret.
		Password: "",

		// Database name (schema) to connect to.
		// Each environment should have a separate database.
		// Database must exist before application startup.
		// Oracle: this is the service name (e.g. "FREEPDB1"), not the SID.
		// Default: "gorest"; if database name not configure
		Name: "gorest",

		// Path to Oracle Instant Client library directory (godror/OCI only).
		// Only needed if Instant Client is not in system default library path.
		// Example: "/opt/oracle/instantclient_21_8"
		// Default: empty (assumes Instant Client is in system PATH/LD_LIBRARY_PATH)
		LibDir: "",

		// TLS/SSL_MODE
		// Transport Layer Security (encryption) for database connection.
		// PostgreSQL: "disable", "require", "verify-ca", "verify-full"
		// MSSQL/SQLSERVER: disable, optional, mandatory, strict
		// SQLite: leave empty (not applicable)
		// MongoDB: true, false (default: false)
		// Oracle: "enable", "disable" — controls whether connection uses TCPS (TLS)
		// protocol instead of plain TCP. When "enable", a wallet (see SSLCAPath)
		// is required to establish the encrypted connection. Unlike Postgres/MSSQL,
		// Oracle does not support granular modes like "verify-ca"; it's effectively on/off.
		// Development: use "disable" (faster for localhost)
		// Production: use "verify-full" or "strict" (maximum security); Oracle: use "enable"
		// Default: "disable" (change this explicitly in production!); if driver configured using mongo, default to "false"
		TLS: "disable",

		// Path to SSL certificate file (client certificate for mutual TLS).
		// Example: "./data/cert/client.crt"
		// Not applicable for Oracle: godror/OCI uses wallet-based auth (see SSLCAPath),
		// not separate client cert/key files like Postgres/MySQL mTLS. Leave empty for Oracle.
		// Default: empty
		SSLCertPath: "",

		// Path to SSL key file (client private key).
		// Example: "./data/cert/client.key"
		// Not applicable for Oracle: same reason as SSLCertPath, Oracle wallet bundles
		// certificate/key/trust material together instead of using discrete files. Leave empty for Oracle.
		// Default: empty
		SSLKeyPath: "",

		// Path to CA certificate file for server certificate verification.
		// Example: "./data/cert/ca.crt"
		// Oracle: repurposed as the path to the Oracle Wallet directory (equivalent to
		// TNS_ADMIN), containing cwallet.sso, tnsnames.ora, sqlnet.ora, etc. Required
		// when TLS is "enable". Example: "./data/oracle/wallet"
		// Default: empty
		SSLCAPath: "",

		// Character set for connection (MySQL).
		// Not applicable for Oracle: character set is controlled via the NLS_LANG
		// environment variable (e.g. "AMERICAN_AMERICA.AL32UTF8") on the client machine,
		// not passed through the connection string. Leave empty for Oracle.
		// Default: "utf8mb4" (recommended)
		Charset: "utf8mb4",

		// Path to database files directory.
		// Example: "./database"
		// Default: "./database"
		DatabaseDirectory: "./app/database",

		// Migration directory, relative to DatabaseDirectory.
		// The driver name is appended automatically.
		// Example: "migrations" -> ./app/database/migrations/{Driver}
		// Default: "migrations"
		MigrationDirectory: "migrations",

		// Seeder directory, relative to DatabaseDirectory.
		// The driver name is appended automatically.
		// Example: "seeders" -> ./app/database/seeders/{Driver}
		// Default: "seeders"
		SeederDirectory: "seeders",

		// Timeout duration limit for database connection and operations.
		// Applied to DialTimeout, ReadTimeout, and WriteTimeout.
		// Must be greater than or equal to 0, Please do not set < 0 or lower then 0 (minus)
		// Recommended: 5 seconds (for production)
		// Default: 5 seconds
		Timeout: 5 * time.Second,

		// Maximum lifetime of connection before being closed and recreated.
		// Prevents stale connections and memory leaks.
		// Must be greater than or equal to 0, Please do not set < 0 or lower then 0 (minus)
		// Recommended: 10-30 minutes
		// Default: 20 minutes
		ConnectionMaxLifetime: 20 * time.Minute,

		// MSSQL/SQLSERVER only: whether to trust self-signed server certificates.
		// false: verify certificate (production, secure)
		// true: trust any certificate (development/testing only - NEVER true in production!)
		// Default: true
		TrustServerCertificate: true,

		// Enable query logging/debugging.
		// true: log all queries (slow, development only)
		// false: no logging (production)
		// Default: false
		QueryLogging: false,

		// Whether to auto-run database migrations on application startup.
		// Development: true (fast iteration)
		// Production: false (manual control, recommended)
		// Default: false
		Migration: true,

		// Whether to auto-run database seeders on application startup.
		// Development: true (fast iteration)
		// Production: false (manual control, recommended)
		// Default: false
		Seeder: true,

		// Maximum number of open connections to database.
		// Default: unlimited (change this explicitly in production!)
		// Must be greater than or equal to 0, Please do not set < 0 or lower then 0 (minus)
		// High-traffic recommended: 50-100
		// Default: 50
		MaxOpenConnections: 50,

		// Maximum number of idle connections in the pool.
		// Must be greater than or equal to 0, Please do not set < 0 or lower then 0 (minus)
		// Default: 5
		MaxIdleConnections: 5,
	}
}
