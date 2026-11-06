# PostgreSQL configuration

This starter's application data layer uses PostgreSQL. Runtime settings are read by `app/config/database.go` and mapped to GOREST's database configuration.

## Application connection

Set these variables in `.env`:

| Variable | Purpose |
|---|---|
| `GOREST_DB_DRIVER` | Must be `postgres` for this starter. |
| `GOREST_DB_HOST`, `GOREST_DB_PORT` | PostgreSQL host/IP and port (default `5432`). |
| `GOREST_DB_NAME` | Application database name. |
| `GOREST_DB_USER`, `GOREST_DB_PASSWORD` | Credentials; both are required. |
| `GOREST_DB_TLS` | PostgreSQL SSL mode: prefer `verify-full` to verify the server certificate and hostname. |
| `GOREST_DB_CA_PATH` | Optional PEM trust certificate for a private server CA. |
| `GOREST_DB_CERT_PATH`, `GOREST_DB_KEY_PATH` | Optional client certificate and matching private key for mTLS; configure both together. |
| `GOREST_DB_MIGRATION` | Enables the PostgreSQL migration during application startup. |
| `GOREST_DB_SEEDER` | Enables sample data insertion; disabled by default. |

The `users` table is created by `app/database/migrations/postgres/001_create_users_table.sql`. Set `GOREST_DB_MIGRATION=true` for a controlled migration run, then turn it off for normal application startup. The down migration drops the `users` table; never run it against data that must be retained. The sample seeder inserts ten rows with placeholder password hashes for demonstrating list responses; these are not usable login credentials. Enable it only in a disposable or dedicated non-production database.

Use a dedicated database account with only the required privileges. Do not copy production credentials into test files, shell history, or source control. Integration tests require a dedicated non-production database and use independent safeguards before running migrations, seeders, or test data writes.

## Configure TLS, private CAs, and mutual TLS

Set `GOREST_DB_TLS=verify-full` to verify the PostgreSQL server certificate and hostname. `verify-ca` validates the certificate chain but does not check the hostname. `require` encrypts the connection but does not provide equivalent server identity verification, so use it only when that tradeoff is intentional. These are the TLS modes accepted by this starter's PostgreSQL integration tests. Ensure `GOREST_DB_HOST` matches a DNS name or IP address in the server certificate's subject alternative names.

When the server certificate chains to a CA trusted by the host system, a custom CA path is not needed. For a private CA, set `GOREST_DB_CA_PATH` to a readable PEM file containing the trust chain required to verify the PostgreSQL server.

For mTLS, configure `GOREST_DB_TLS=verify-full`, the server-verifying CA in `GOREST_DB_CA_PATH`, and both `GOREST_DB_CERT_PATH` and `GOREST_DB_KEY_PATH` for the client certificate and matching private key. The PostgreSQL server must trust the client certificate issuer and be configured to require client-certificate authentication. The server-side `pg_hba.conf` policy and certificate trust configuration determine whether a client certificate is required; client-side settings alone do not enable server enforcement.

Keep private keys readable only by the application account, and never disable certificate verification in production. See the [PostgreSQL TLS documentation](https://www.postgresql.org/docs/current/libpq-ssl.html), [GOREST v1.1.20 configuration guide](https://github.com/elsyahtech/gorest/blob/v1.1.20/docs/configuration.md), and [testing guide](./testing.md) for client-side setup and integration-test safeguards.

## Integration-test database

The PostgreSQL integration suite uses the `GOREST_DB_*` connection variables and requires `GOREST_TEST_POSTGRES_CONFIRM_DATABASE` to match `GOREST_DB_NAME`. Database names ending in `_test` or `_testing` are recognized as test databases. For another dedicated non-production name, set `GOREST_TEST_POSTGRES_ALLOW_NON_TEST_DATABASE` to the exact same name as the confirmation variable. The HTTP integration test runs the PostgreSQL migration and seeder, then inserts and removes only its UUID-scoped test user. Never use these overrides or tests against production.
