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
| `GOREST_DB_CA_PATH` | Not required for MySQL TLS; GOREST v1.1.19 does not use this field for MySQL. |
| `GOREST_DB_CERT_PATH`, `GOREST_DB_KEY_PATH` | Present in the shared GOREST config surface but not applied to the MySQL client's TLS profile. |
| `GOREST_DB_MIGRATION` | Enables the MySQL migration during application startup. |
| `GOREST_DB_SEEDER` | Enables the sample seeder; disabled in the example configuration. |

The `users` table is created by `app/database/migrations/mysql/001_create_users_table.sql`. Set `GOREST_DB_MIGRATION=true` for a controlled migration run, then turn it off for normal application startup. The down migration drops the `users` table; never run it against data that must be retained.

Use a dedicated database account with only the required privileges. Do not copy production credentials into test files, shell history, or source control. The application does not make an application-level claim that an arbitrary `GOREST_DB_NAME` is safe for tests; the integration test has its own independent guardrails.

## Configure TLS, private CAs, and client certificates

Set `GOREST_DB_TLS=true` to require MySQL TLS. GOREST v1.1.19 passes the driver's built-in `tls=true` option, which uses the Go/system trusted roots and verifies the server name against `GOREST_DB_HOST`. A CA path is not required when the server certificate chains to a trusted system root.

Ensure `GOREST_DB_HOST` matches a DNS name or IP address in the server certificate's subject alternative names.

For a private CA, add the CA PEM used to verify the server (the same trust certificate supplied to a successful `mysql --ssl-ca` connection) to the trusted roots available to the application process before startup. Do not substitute the server's leaf certificate or a file merely named `fullchain` unless your database administrator confirms it is the CA trust bundle. On Linux, create a bundle containing the system CA bundle plus the private CA certificate and set `SSL_CERT_FILE` to that bundle when starting the application. For example:

```sh
cat /etc/ssl/certs/ca-certificates.crt /path/to/private-ca.pem > /tmp/mysql-trust-bundle.pem
chmod 600 /tmp/mysql-trust-bundle.pem
SSL_CERT_FILE=/tmp/mysql-trust-bundle.pem GOREST_DB_TLS=true go run .
```

On macOS or Windows, add the CA to the operating-system trust store used by Go, or use a platform-supported trusted-root bundle and verify that the application process can read it. Do not disable certificate verification or use a `skip-verify` TLS mode.

The current GOREST MySQL connector does not load `GOREST_DB_CERT_PATH` and `GOREST_DB_KEY_PATH` into a MySQL client TLS profile, so this starter does not provide mTLS configuration through those variables. If your server requires client-certificate authentication, use a connector implementation that explicitly supports loading the CA and client certificate/key. Consult the [GOREST documentation and source](https://github.com/elsyahtech/gorest/tree/main), the [go-sql-driver/mysql TLS configuration](https://github.com/go-sql-driver/mysql#tls), and the official [MySQL encrypted-connection documentation](https://dev.mysql.com/doc/refman/8.4/en/using-encrypted-connections.html). Validate both successful and rejected certificate cases before deployment; a non-empty `Ssl_cipher` confirms encryption, but not which trust roots were used.

## Integration-test database

Integration tests use the `GOREST_TEST_MYSQL_*` variables documented in the single `.env.example` template. They take precedence over application connection settings. Missing test connection values can fall back to `GOREST_DB_*` only for a MySQL database in a development or testing environment. The test harness normally requires a database name ending in `_test` or `_testing` and requires `GOREST_TEST_MYSQL_CONFIRM_DATABASE` to match it exactly before it writes test data.

The tests do not create or drop tables. Apply the MySQL migration to the test schema before testing; tests create a unique user via the HTTP API, verify the row directly, and delete only that unique ID/email pair during cleanup. See [Testing and coverage](./testing.md) for environment setup and report commands.
