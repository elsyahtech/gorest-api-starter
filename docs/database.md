# SQL Server configuration

This starter's application data layer uses Microsoft SQL Server. Runtime settings are read by `app/config/database.go` and mapped to GOREST's database configuration.

## Application connection

Set these variables in `.env`:

| Variable | Purpose |
|---|---|
| `GOREST_DB_DRIVER` | Must be `sqlserver` for this starter. |
| `GOREST_DB_HOST`, `GOREST_DB_PORT` | SQL Server hostname/IP and TCP port (default `1433`). |
| `GOREST_DB_NAME` | Existing application database name. |
| `GOREST_DB_USER`, `GOREST_DB_PASSWORD` | SQL Server login credentials; both are required. |
| `GOREST_DB_TLS` | SQL Server encryption mode; `mandatory` is the default, and `strict` is also supported. |
| `GOREST_DB_CA_PATH` | Optional PEM trust bundle for a private SQL Server CA. |
| `GOREST_DB_TRUST_SERVER_CERTIFICATE` | Defaults to `false`; leave false to validate the SQL Server certificate and hostname. |
| `GOREST_DB_MIGRATION` | Enables the SQL Server migration during application startup. |
| `GOREST_DB_SEEDER` | Enables sample data insertion; disabled by default. |

The database must exist before application startup. The `users` table is created by `app/database/migrations/sqlserver/001_create_users_table.sql`. Set `GOREST_DB_MIGRATION=true` for a controlled migration run, then turn it off for normal application startup. The down migration drops the `users` table; never run it against data that must be retained. The sample seeder inserts ten rows with placeholder password hashes for demonstrating list responses; these are not usable login credentials. Enable it only in a disposable or dedicated non-production database.

Use a dedicated database account with only the required privileges. Applying migrations and seeders requires permission to create the migration/seeder tracking tables, create the users table and index, and insert sample users. Do not copy production credentials into test files, shell history, or source control.

## TLS and server certificate validation

SQL Server connections require encryption by default through `GOREST_DB_TLS=mandatory`. `strict` is also supported. With `GOREST_DB_TRUST_SERVER_CERTIFICATE=false`, the SQL Server driver validates the server certificate chain against the host's trusted roots and checks the server name. The hostname in `GOREST_DB_HOST` must match a DNS name or IP address in the certificate's subject alternative names. For a private CA, set `GOREST_DB_CA_PATH` to a PEM trust bundle; it is passed to the SQL Server driver for chain validation.

Keep `GOREST_DB_TRUST_SERVER_CERTIFICATE=false` in production. Setting it to `true` accepts an untrusted server certificate and weakens protection against interception; the integration test runner rejects this setting. Do not disable TLS for remote connections.

See the [Microsoft SQL Server TLS documentation](https://learn.microsoft.com/en-us/sql/database-engine/configure-windows/tls-1-2-support), the [GOREST v1.1.21 configuration guide](https://github.com/elsyahtech/gorest/blob/v1.1.21/docs/configuration.md), and the [testing guide](./testing.md) for connection and integration-test setup.

## Integration-test database

The SQL Server integration suite uses the `GOREST_DB_*` connection variables and requires `GOREST_TEST_SQLSERVER_CONFIRM_DATABASE` to match `GOREST_DB_NAME` exactly. It refuses to run against SQL Server system databases or with `GOREST_APP_ENV=production`. The HTTP integration test runs the SQL Server migration and seeder, then inserts and removes only its UUID-scoped test user. It does not truncate tables or drop the database. Use a dedicated, non-production database and never use the test confirmation against production.
