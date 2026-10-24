# Database, migrations, and seeders

The application database settings are loaded in `app/config/database.go`. The starter begins with `database.ConfigDefault`, which selects SQLite, and enables this starter's migrations and seeders (`Migration=true`, `Seeder=true`) for local examples. The local examples need no external database service or credentials. These flags are **local-development settings, not production settings**. With `APP_ENV=production`, migration and seeder defaults are false and startup rejects attempts to enable them. The application reads environment variables directly, but does not load `.env` files; inject secrets through your shell, container/orchestrator, or secret manager.

## Use another database driver

For a server database, set the environment variables below. The application applies these values at startup:

```sh
export APP_ENV=development
export GOREST_DB_HOST=127.0.0.1
export GOREST_DB_NAME=gorest_dev
export GOREST_DB_USERNAME=gorest_app
export GOREST_DB_PASSWORD='<local-development-password>'
export GOREST_DB_MIGRATION=true
export GOREST_DB_SEEDER=true
```

Replace the sample host, database, username, and password with values for your disposable development database. The password line is a placeholder; inject the real value from your shell or secret manager. Then select one driver configuration below. These examples are for local development; verify the database TLS settings before using a remote or production service.

### PostgreSQL

```sh
export GOREST_DB_DRIVER=postgres
export GOREST_DB_PORT=5432
export GOREST_DB_TLS=disable
```

### MySQL

```sh
export GOREST_DB_DRIVER=mysql
export GOREST_DB_PORT=3306
export GOREST_DB_TLS=disable
```

For TLS, GOREST passes the driver TLS mode when `GOREST_DB_SSL_CA_PATH` is non-empty. `GOREST_DB_TLS=true` uses the driver's system certificate roots; a custom CA or client certificate requires a registered MySQL TLS profile. GOREST does not load a custom CA from `GOREST_DB_SSL_CA_PATH` for MySQL.

### Oracle

```sh
export GOREST_DB_DRIVER=oracle
export GOREST_DB_PORT=1521
export GOREST_DB_TLS=disable
# Optional when Oracle Instant Client is not available through the system library path:
export GOREST_DB_LIB_DIR=/path/to/instantclient
```

Set `GOREST_DB_NAME` to the Oracle **service name** (for example, `FREEPDB1`), not the SID. For TCPS, set `GOREST_DB_TLS=enable` and `GOREST_DB_SSL_CA_PATH` to the Oracle wallet/configuration directory. The selected service must have Oracle Instant Client available.

### SQL Server

```sh
export GOREST_DB_DRIVER=sqlserver
export GOREST_DB_PORT=1433
export GOREST_DB_TLS=disable
```

For encrypted connections, use `GOREST_DB_TLS=mandatory` (or `strict` when supported by the server and driver). The starter leaves `TrustServerCertificate` false, so the server certificate is validated; provide a certificate trusted by the client. The current environment loader does not expose `TrustServerCertificate` as an override.

### MongoDB

```sh
export GOREST_DB_DRIVER=mongo
export GOREST_DB_PORT=27017
export GOREST_DB_TLS=false
# Optional: set both values when authentication is enabled.
export GOREST_DB_USERNAME=gorest_app
export GOREST_DB_PASSWORD='<local-development-password>'
```

`GOREST_DB_NAME` is the database name. When credentials are set, this GOREST version uses `authSource=admin`; unauthenticated connections omit credentials. For TLS, set `GOREST_DB_TLS=true`, and use `GOREST_DB_SSL_CA_PATH` for a CA file or `GOREST_DB_SSL_CERT_PATH` for a client certificate/key PEM file when needed.

### ScyllaDB

```sh
export GOREST_DB_DRIVER=scylla
export GOREST_DB_PORT=9042
# GOREST_DB_NAME is the keyspace; create it before connecting.
export GOREST_DB_NAME=gorest_dev
# Optional: set both values when password authentication is enabled.
export GOREST_DB_USERNAME=gorest_app
export GOREST_DB_PASSWORD='<local-development-password>'
```

The current GOREST database configuration does not expose ScyllaDB TLS options through `GOREST_DB_TLS` or the other `GOREST_DB_*` variables. For a remote ScyllaDB service, configure TLS through a driver/network layer supported by your deployment, or use a trusted encrypted network. In production, GOREST warns that this database setting does not enable ScyllaDB transport encryption; the warning does not prevent startup.

The default ports shown above are common service defaults; use the actual port and endpoint from your environment. MySQL, PostgreSQL, SQL Server, and Oracle expect a database/schema or service name; MongoDB expects a database name; ScyllaDB expects an existing keyspace. Migration and seeder files are driver-specific. Review them before enabling the local `true/true` example settings against a database with existing data.

For production, set `APP_ENV=production` and inject the database connection fields. The starter requires `GOREST_DB_TLS` to be explicitly set for MySQL, PostgreSQL, Oracle, SQL Server, and MongoDB, but it does **not** require that the selected mode enables encryption: for example, `disable` is accepted where that is a supported driver mode. ScyllaDB does not use this setting, and SQLite has no network transport to encrypt. This explicit setting requirement makes the selected behavior visible; it is not a universal TLS enforcement policy. Review the table below, GOREST's startup warning, and your own network/security policy before choosing a mode. For SQLite, set `GOREST_DB_DRIVER=sqlite` and an explicit persistent `GOREST_DB_DIRECTORY`. Set `GOREST_DB_MIGRATION=false` and `GOREST_DB_SEEDER=false`; these are production defaults, and startup fails if either is enabled. Use a secret manager or protected deployment environment for passwords. The application does not load `.env` files.

### Production transport encryption by driver

The starter provides connection settings and GOREST emits a warning in production when transport encryption is disabled or cannot be confirmed from the configuration. That warning is informational: it does not block startup and does not configure TLS for you. No single TLS setting works for every driver. Choose a mode compatible with your database and deployment; if a client, server, or network path cannot use TLS, decide whether a separately protected network meets your policy. A private IP address or firewall alone does not encrypt traffic.

| Driver | TLS behavior available through this starter | Production considerations |
| :-- | :-- | :-- |
| SQLite | Not applicable; the database is a local file, not a network connection. | Protect the file, its directory, and backups with host/storage access controls. |
| MySQL | `GOREST_DB_TLS=true` uses the driver's system certificate roots. GOREST passes the TLS option only when `GOREST_DB_SSL_CA_PATH` is non-empty; this path is a switch and is not loaded as a custom CA. | The default `disable` mode is accepted but triggers a warning. For a custom CA or client certificate, the Go MySQL driver needs a registered TLS profile; this starter does not expose profile registration as an environment setting. |
| PostgreSQL | `GOREST_DB_TLS` maps to `sslmode`: `disable`, `require`, `verify-ca`, or `verify-full`. CA, client certificate, and key paths use the corresponding `GOREST_DB_SSL_*_PATH` settings. | `require` encrypts without verifying the server identity; `verify-ca` checks the CA chain but not the hostname. `verify-full` with a trusted CA provides certificate and hostname verification. `allow` and `prefer` can fall back to plaintext. |
| Oracle | `enable` selects TCPS and uses `GOREST_DB_SSL_CA_PATH` as the Oracle wallet/configuration directory; `disable` uses TCP. | TCPS needs a correctly configured wallet and compatible Oracle service. `GOREST_DB_LIB_DIR` selects the Instant Client location when it is not on the system library path. |
| SQL Server | `GOREST_DB_TLS` maps to the driver's `encrypt` setting (for example, `mandatory` or `strict`). The starter keeps `TrustServerCertificate=false`. | Use a server certificate trusted by the client. This starter does not expose `TrustServerCertificate` as an environment override. |
| MongoDB | `GOREST_DB_TLS=true` enables TLS; `false` disables it. CA and client certificate/key PEM paths can be supplied through `GOREST_DB_SSL_CA_PATH` and `GOREST_DB_SSL_CERT_PATH`. | Confirm that the server certificate is trusted and that the configured server topology supports the selected TLS setup. |
| ScyllaDB | `GOREST_DB_TLS` and the starter's other database TLS environment variables do not configure ScyllaDB TLS. | GOREST warns in production that transport encryption is not configured by this field. Configure TLS using ScyllaDB driver options outside this starter, or protect the path with a trusted encrypted network. |

The warning is a useful signal, not a security guarantee: it cannot assess every proxy, tunnel, certificate deployment, or network control. Conversely, the absence of a warning does not certify the entire deployment. The application owner chooses and verifies the policy appropriate to the service and environment.

The examples use GOREST's driver names: `mysql`, `postgres`, `oracle`, `sqlserver`, `mongo`, and `scylla`. Review the selected driver's migration and seeder files, and adapt native queries because SQL syntax, constraints, and transaction support differ. See the [GOREST v1.0.2 documentation](https://github.com/elsyahtech/gorest/tree/v1.0.2/docs) for framework-level details.

## Drivers and file locations

GOREST supports MySQL, PostgreSQL, SQL Server, Oracle, SQLite, MongoDB, and ScyllaDB. This starter defaults to SQLite. Migration and seeder files are grouped by driver:

```text
app/database/
  migrations/<driver>/
  seeders/<driver>/
```

Migrations define the `departments`, `users`, and `orders` tables or collections. Seeders provide sample data for those resources. The PostgreSQL examples include departments with IDs `dept-01` through `dept-05`, sample users, and sample orders. Check the files for your selected driver before relying on specific IDs or data.

The maintainer reports manually exercising all seven database integrations successfully. The automated `make test-e2e` runner provisions disposable services and runs migration, seeder, CRUD, applicable transaction, lifecycle, concurrency, and native-query route scenarios for the supported drivers. See the [driver verification matrix](./testing.md#driver-verification-matrix) for scenario applicability and recorded evidence.

## Paths and startup settings

`app/config/database.go` inherits connection defaults such as the SQLite driver, local database name, and empty credentials from `database.ConfigDefault`. It overrides these starter-specific settings:

- `DatabaseDirectory` points to `./app/database`.
- `Migration` and `Seeder` are enabled so the matching SQLite files run when the application starts. Bare GOREST `ConfigDefault` leaves migrations and seeders disabled.
- The SQLite example uses a single connection in its pool so simultaneous HTTP writes queue against SQLite's single-writer model. Server database drivers inherit GOREST's connection-pool defaults.

When selecting a server driver, configure its host, port, username, password, database/service name, and any driver-supported TLS mode and certificate or wallet paths. GOREST's defaults for timeouts and connection pool limits remain in effect unless you override them. See [Production transport encryption by driver](#production-transport-encryption-by-driver) for what the starter does and does not configure.

Directory paths are relative to the process working directory. Run the application from the repository root as described in [Getting started](./getting-started.md), or adjust the paths if your deployment layout is different.

### Production startup

The starter's local `true/true` settings must not be used unchanged in production. Set `Migration=false` and `Seeder=false` for normal production application instances. Apply schema changes once through a controlled deployment migration job before serving traffic. Do not run these sample seeders against production data; if a production system needs initial/reference data, provide a deliberate, reviewed, and repeatable deployment process for that data instead.

GOREST runs migrations before seeders and does not start the application if a startup migration or seeder fails. A failed script may have already changed the database, so inspect the database and migration history before retrying it. Treat applied migration and seeder files as immutable; add a new file for later changes instead of editing an applied file. For the framework's one-time database startup task workflow and recovery details, see the [GOREST v1.0.2 migration and seeder guide](https://github.com/elsyahtech/gorest/blob/v1.0.2/docs/guide.md).

## Example relationships

- Users store `department_id`; the SQL schema also creates an index for this column.
- Orders store `user_id`; the PostgreSQL schema defines a foreign key from orders to users.
- The PostgreSQL SQL schema requires user email addresses to be unique.

Check the migration files for your selected driver before assuming that constraints or behavior are the same across all databases.

## Native queries and transactions

The native SQL user routes select parameter markers for the configured SQL driver and are tested against SQLite, PostgreSQL, MySQL, Oracle, and SQL Server. SQL syntax and behavior may still vary between drivers. Use the MongoDB routes only with MongoDB, and the Scylla routes only with ScyllaDB.

Some ORM operations use `BeginTx`, `Commit`, and `Rollback`. In these examples, GOREST transaction calls are no-ops for MongoDB and ScyllaDB. Do not assume that writes using those drivers are atomic transactions.

For native route details, see the [API reference](./api-reference.md). For driver settings and ORM APIs, see the [GOREST v1.0.2 documentation](https://github.com/elsyahtech/gorest/tree/v1.0.2/docs).
