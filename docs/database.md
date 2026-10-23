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

The current GOREST database configuration does not expose ScyllaDB TLS options. Configure TLS at a supported driver/network layer before using a remote ScyllaDB service.

The default ports shown above are common service defaults; use the actual port and endpoint from your environment. MySQL, PostgreSQL, SQL Server, and Oracle expect a database/schema or service name; MongoDB expects a database name; ScyllaDB expects an existing keyspace. Migration and seeder files are driver-specific. Review them before enabling the local `true/true` example settings against a database with existing data.

For production, set `APP_ENV=production` and inject `GOREST_DB_DRIVER`, `GOREST_DB_HOST`, `GOREST_DB_NAME`, `GOREST_DB_USERNAME`, `GOREST_DB_PASSWORD`, and `GOREST_DB_TLS` explicitly for server databases. For SQLite, set `GOREST_DB_DRIVER=sqlite` and an explicit persistent `GOREST_DB_DIRECTORY`. Set `GOREST_DB_MIGRATION=false` and `GOREST_DB_SEEDER=false`; these are production defaults, and startup fails if either is enabled. Use a secret manager or protected deployment environment for passwords. The application does not load `.env` files.

The examples use GOREST's driver names: `mysql`, `postgres`, `oracle`, `sqlserver`, `mongo`, and `scylla`. Review the selected driver's migration and seeder files, and adapt native queries because SQL syntax, constraints, and transaction support differ. See the [GOREST v1.0.2 documentation](https://github.com/elsyahtech/gorest/tree/v1.0.2/docs) for framework-level details.

## Drivers and file locations

GOREST supports MySQL, PostgreSQL, SQL Server, Oracle, SQLite, MongoDB, and ScyllaDB. This starter defaults to SQLite. Migration and seeder files are grouped by driver:

```text
app/database/
  migrations/<driver>/
  seeders/<driver>/
```

Migrations define the `departments`, `users`, and `orders` tables or collections. Seeders provide sample data for those resources. The PostgreSQL examples include departments with IDs `dept-01` through `dept-05`, sample users, and sample orders. Check the files for your selected driver before relying on specific IDs or data.

The maintainer reports manually exercising all seven database integrations successfully. Separately, the `test/drivers` integration checks have automated migration, seeder, CRUD, and applicable transaction verification for all seven drivers. Server-driver checks require the corresponding service and driver-specific connection settings. The HTTP E2E suite primarily uses SQLite, with opt-in native MongoDB and ScyllaDB route tests. See the [driver verification matrix](./testing.md#driver-verification-matrix) for the per-operation results and limitations.

## Paths and startup settings

`app/config/database.go` inherits connection defaults such as the SQLite driver, local database name, and empty credentials from `database.ConfigDefault`. It overrides these starter-specific settings:

- `DatabaseDirectory` points to `./app/database`.
- `Migration` and `Seeder` are enabled so the matching SQLite files run when the application starts. Bare GOREST `ConfigDefault` leaves migrations and seeders disabled.
- The SQLite example uses a single connection in its pool so simultaneous HTTP writes queue against SQLite's single-writer model. Server database drivers inherit GOREST's connection-pool defaults.

When selecting a server driver, configure its host, port, username, password, database/service name, TLS, and any certificate or wallet paths as required by that driver. GOREST's defaults for timeouts and connection pool limits remain in effect unless you override them.

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

The native SQL user routes use PostgreSQL syntax and `$1` placeholders. Adapt the query and placeholders before using a different SQL driver. Use the MongoDB routes only with MongoDB, and the Scylla routes only with ScyllaDB.

Some ORM operations use `BeginTx`, `Commit`, and `Rollback`. In these examples, GOREST transaction calls are no-ops for MongoDB and ScyllaDB. Do not assume that writes using those drivers are atomic transactions.

For native route details, see the [API reference](./api-reference.md). For driver settings and ORM APIs, see the [GOREST v1.0.2 documentation](https://github.com/elsyahtech/gorest/tree/v1.0.2/docs).
