# Testing your PostgreSQL starter

The starter provides fast unit tests and optional integration tests against a PostgreSQL server you control. Unit tests do not connect to a database. Integration tests exercise the real database connection, including mTLS, library operations, startup migration/seeding, and HTTP endpoints.

## Run unit tests

Unit tests are useful for checking isolated configuration, validation, routing, and response-mapping behavior:

```sh
make unit-test
# equivalent:
go test ./... -count=1
```

They need no database credentials or external services. They do not prove that PostgreSQL connectivity, SQL statements, migrations, or TLS settings work against your server.

## Prepare a PostgreSQL integration database

Integration tests use `GOREST_DB_*` connection settings from the process environment or `.env`. Before running them:

1. Use a dedicated, non-production PostgreSQL database and a restricted database account. The HTTP flow runs the configured migration and seeder, so the account needs the permissions required for those operations and the tested user operations.
2. Set `GOREST_DB_DRIVER=postgres`, `GOREST_DB_DISABLE=false`, and the host, port, database name, username, and password.
3. Configure TLS. `verify-full` is recommended because it verifies both the server certificate chain and hostname. `verify-ca` verifies the chain without hostname checking; `require` encrypts the connection but does not provide equivalent server identity verification.
4. Configure `GOREST_DB_CERT_PATH` and `GOREST_DB_KEY_PATH` with a matching client certificate/key pair. Set `GOREST_DB_CA_PATH` when the PostgreSQL server uses a private CA.
5. Confirm the exact database name by setting `GOREST_TEST_POSTGRES_CONFIRM_DATABASE` to `GOREST_DB_NAME`. Names ending in `_test` or `_testing` are preferred. If your dedicated non-production database has another name, also set `GOREST_TEST_POSTGRES_ALLOW_NON_TEST_DATABASE` to that exact same name.

The confirmation is required even for databases with a `_test` or `_testing` suffix. The test runner refuses production mode and an incomplete client certificate/key pair. Do not point these settings at production.

## Run integration tests

Run the PostgreSQL integration tests:

```sh
make test-integration
```

The suite checks that the server rejects a connection without the client certificate and accepts an encrypted mTLS connection with it. It also calls the user library directly against PostgreSQL, checking user creation and password hashing, listing, and rejection of an incorrect password. The HTTP flow starts the application with migrations and seed data enabled, then checks user creation, duplicate rejection, login, and authenticated listing. Test-created users are uniquely identified and removed afterward; the suite does not truncate tables or drop the database.

To run only the PostgreSQL mTLS connection check:

```sh
make test-postgres-mtls
```

The `integration` build tag keeps external-service checks out of the default unit-test command. The target above selects PostgreSQL checks only; it does not run optional Redis tests.

## Optional Redis checks

Redis is not required by the user API and is disabled by default. If you have enabled Redis for your application, the full tagged suite can also test it when you run:

```sh
make integration-reports
```

That command runs the tagged integration suite once and creates coverage and test-output files under `coverage/`. Redis tests use the configured Redis endpoint when enabled, so point it only at a dedicated non-production instance. They use a unique key and do not flush the database. Redis mTLS checks require an endpoint configured to require client certificates. See [Redis configuration](./configuration.md#redis-optional) before enabling Redis.

## Automated checks in GitHub Actions

The workflow runs unit tests and coverage on pushes and pull requests. Its PostgreSQL integration job runs only when configured with secrets for a reachable, dedicated non-production mTLS database:

- `GOREST_TEST_POSTGRES_HOST`, optional `GOREST_TEST_POSTGRES_PORT` (defaults to `5432`), `GOREST_TEST_POSTGRES_DATABASE`, `GOREST_TEST_POSTGRES_USER`, and `GOREST_TEST_POSTGRES_PASSWORD`.
- `GOREST_TEST_POSTGRES_TLS`, set to `verify-full` or `verify-ca`.
- `GOREST_TEST_POSTGRES_CA_PEM`, `GOREST_TEST_POSTGRES_CLIENT_CERT_PEM`, and `GOREST_TEST_POSTGRES_CLIENT_KEY_PEM`.
- `GOREST_TEST_POSTGRES_CONFIRM_DATABASE`, exactly matching the database name. For a database name not ending in `_test` or `_testing`, also set `GOREST_TEST_POSTGRES_ALLOW_NON_TEST_DATABASE` to that same name.

The workflow stores certificate material in restricted-permission temporary files. Configure the PostgreSQL server to trust and require the client certificate. If the required secrets are unavailable, GitHub skips the integration job; a skipped job is not a passing test result. Credentials and private keys belong in Actions secrets or a secret manager, never in the workflow file.

## Coverage

Generate unit coverage locally with:

```sh
make coverage
```

Generate integration test output and application coverage with:

```sh
make integration-reports
```

These reports help you see which code paths your tests exercised; they do not replace reviewing the behavior and test cases relevant to your application. Regenerate them after source changes if you use them to evaluate coverage.
