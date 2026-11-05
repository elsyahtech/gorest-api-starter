# MySQL configuration

This branch's application data layer is configured for MySQL. The runtime variables are read by `app/config/database.go` and mapped to GOREST's database configuration.

## Application connection

Set these variables in `.env`:

| Variable | Purpose |
|---|---|
| `GOREST_DB_DRIVER` | Must be `mysql` for this starter. |
| `GOREST_DB_HOST`, `GOREST_DB_PORT` | MySQL host/IP and port. |
| `GOREST_DB_NAME` | Application database/schema name. |
| `GOREST_DB_USER`, `GOREST_DB_PASSWORD` | Credentials; both are required. |
| `GOREST_DB_TLS` | Set to `true` to require TLS and verify the server against trusted roots. |
| `GOREST_DB_CA_PATH` | Optional PEM trust certificate for a private server CA. |
| `GOREST_DB_CERT_PATH`, `GOREST_DB_KEY_PATH` | Optional client certificate and private key for MySQL mTLS; configure both together. |
| `GOREST_DB_MIGRATION` | Enables the MySQL migration during application startup. |
| `GOREST_DB_SEEDER` | Enables the sample seeder; disabled in the example configuration. |

The `users` table is created by `app/database/migrations/mysql/001_create_users_table.sql`. Set `GOREST_DB_MIGRATION=true` for a controlled migration run, then turn it off for normal application startup. The down migration drops the `users` table; never run it against data that must be retained.

Use a dedicated database account with only the required privileges. Do not copy production credentials into test files, shell history, or source control. The application does not make an application-level claim that an arbitrary `GOREST_DB_NAME` is safe for tests; the integration test has its own independent guardrails.

## Configure TLS, private CAs, and client certificates

Set `GOREST_DB_TLS=true` to require MySQL TLS. GOREST v1.1.20 uses the Go/system trusted roots and verifies the server name against `GOREST_DB_HOST` when no CA path is set. A CA path is not required when the server certificate chains to a trusted system root.

Ensure `GOREST_DB_HOST` matches a DNS name or IP address in the server certificate's subject alternative names.

For a private CA, set `GOREST_DB_CA_PATH` to a readable PEM file containing the CA certificate used to verify the server. GOREST loads this certificate as the MySQL TLS trust pool and keeps hostname verification enabled. Do not substitute the server's leaf certificate or a file merely named `fullchain` unless your database administrator confirms it contains the CA trust chain.

For mTLS, set `GOREST_DB_TLS=true`, `GOREST_DB_CA_PATH` to the CA that verifies the server, and both `GOREST_DB_CERT_PATH` and `GOREST_DB_KEY_PATH` to the client certificate and matching private key. The server must trust the client certificate issuer and require client authentication (for example, with a MySQL account configured as `REQUIRE X509`). GOREST v1.1.20 loads and presents this certificate/key pair while verifying the server certificate and hostname. Certificate and key paths are ignored when `GOREST_DB_TLS` names a registered MySQL TLS profile; use a profile only for additional custom TLS settings. Keep private keys readable only by the application account, and never disable certificate verification.

The optional `TestMySQLMutualTLSRequiresClientCertificate` integration test uses the same direct CA and client-certificate path configuration as the application. Configure the `GOREST_TEST_MYSQL_MTLS_*` variables in `.env` or the environment, or allow it to use the runtime `GOREST_DB_*` settings in a development/test environment. Use a dedicated MySQL account that is restricted with `REQUIRE X509`. The test expects the server to reject the account without a client certificate (MySQL error 1045), then verifies a successful connection with the certificate and a non-empty `Ssl_cipher`. See [Testing and coverage](./testing.md#mysql-mutual-tls-integration-test), [GOREST MySQL TLS configuration](https://github.com/elsyahtech/gorest/blob/v1.1.20/docs/configuration.md#mysql), [GOREST MySQL test results](https://github.com/elsyahtech/gorest/blob/v1.1.20/docs/testing-result/testing-result-mysql.md), and [go-sql-driver/mysql TLS configuration](https://github.com/go-sql-driver/mysql#tls). A non-empty `Ssl_cipher` confirms encryption; rejection without the client certificate and success with it verify the endpoint's client-certificate policy.

## Integration-test database

Integration tests use the `GOREST_TEST_MYSQL_*` variables documented in the single `.env.example` template. They take precedence over application connection settings. Missing test connection values can fall back to `GOREST_DB_*` only for a MySQL database in a development or testing environment. The test harness normally requires a database name ending in `_test` or `_testing` and requires `GOREST_TEST_MYSQL_CONFIRM_DATABASE` to match it exactly before it writes test data.

The tests do not create or drop tables. Apply the MySQL migration to the test schema before testing; tests create a unique user via the HTTP API, verify the row directly, and delete only that unique ID/email pair during cleanup. See [Testing and coverage](./testing.md) for environment setup and report commands.
