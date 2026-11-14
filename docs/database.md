# ScyllaDB configuration

This branch's application data layer uses ScyllaDB through GOREST v1.1.21. Runtime variables are read by `app/config/database.go` and mapped to GOREST's database configuration.

## Application connection

Configure these variables in `.env` or the process environment:

| Variable | Purpose |
|---|---|
| `GOREST_DB_DRIVER` | Must be `scylla` for this starter. |
| `GOREST_DB_HOST`, `GOREST_DB_PORT` | ScyllaDB host/IP and native CQL port; the port defaults to `9042`. |
| `GOREST_DB_NAME` | Required name of an already-created keyspace. |
| `GOREST_DB_USER`, `GOREST_DB_PASSWORD` | Optional credentials; configure both or neither. |
| `GOREST_DB_TLS` | `enable` verifies the ScyllaDB server certificate and hostname; `disable` turns TLS off. |
| `GOREST_DB_CA_PATH` | Optional PEM trust certificate for a private server CA. |
| `GOREST_DB_CERT_PATH`, `GOREST_DB_KEY_PATH` | Optional client certificate and private key for mTLS; configure both together. |
| `GOREST_DB_TIMEOUT` | Connection/query timeout, parsed as a Go duration (for example `5s`). |
| `GOREST_DB_MIGRATION` | Enables the ScyllaDB migration during application startup. |
| `GOREST_DB_SEEDER` | Enables sample data; disabled in `.env.example`. |

ScyllaDB does not create the configured keyspace. Provision it before starting the application, then run with `GOREST_DB_MIGRATION=true` to create the users table and email index. Turn migration off for normal startup. The down migration drops the users table and its index; do not run it if the data must be retained. Migration and seeder metadata remain in the keyspace.

Use a dedicated database account with only the permissions required for the selected application and controlled schema-migration operations. Migrations use GOREST's ScyllaDB migration lock and therefore require the corresponding table and lightweight-transaction permissions. Confirm those permissions with your ScyllaDB administrator before running migrations.

## TLS and client certificates

Set `GOREST_DB_TLS=enable` for remote connections. GOREST uses Go's system trust roots when no private CA path is configured and verifies the server name against `GOREST_DB_HOST`. Ensure that the host matches a DNS name or IP address in the server certificate's subject alternative names.

For a private CA, set `GOREST_DB_CA_PATH` to a readable PEM file containing the CA certificate used to verify the server. Do not substitute the server's leaf certificate unless the ScyllaDB administrator explicitly confirms that trust model.

For mTLS, set `GOREST_DB_TLS=enable`, configure the CA path when the server uses a private CA, and set both `GOREST_DB_CERT_PATH` and `GOREST_DB_KEY_PATH` to a matching client certificate and private key. The server must trust the client certificate issuer and require client authentication. The starter rejects an unpaired certificate and key. Restrict key-file permissions and never disable certificate verification for production traffic.

GOREST v1.1.21's ScyllaDB TLS configuration accepts `enable` or `disable`; see the [GOREST configuration guide](https://github.com/elsyahtech/gorest/blob/v1.1.21/docs/configuration.md) and [GOREST ScyllaDB implementation](https://github.com/elsyahtech/gorest/tree/v1.1.21/database).

## Integration-test keyspace

Database-writing tests use the `GOREST_TEST_SCYLLA_*` variables in `.env.example`, taking precedence over application connection settings. Configure an isolated, non-production keyspace and set `GOREST_TEST_SCYLLA_CONFIRM_KEYSPACE` to its exact name before running them. Keyspaces should end in `_test` or `_testing`; another name requires the additional exact-name opt-in `GOREST_TEST_SCYLLA_ALLOW_NON_TEST_KEYSPACE`.

The library and HTTP tests apply the up migration and seeder. The tests do not drop the keyspace, tables, or seed data; they clean up only the uniquely identified user rows they create. Migration/seeder metadata and ten sample rows remain after the run. Use a disposable test keyspace or review and remove those specific test artifacts under your own retention process.

When dedicated test connection settings are absent, `GOREST_DB_*` fallback is allowed only when `GOREST_APP_ENV` is `development`, `test`, or `testing`, the configured driver is ScyllaDB, and the database is enabled. Exact keyspace confirmation is still required. Never point the suite at production.
