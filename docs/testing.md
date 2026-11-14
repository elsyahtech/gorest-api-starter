# Testing your SQL Server starter

The starter provides fast unit tests and optional integration tests against a SQL Server instance you control. Unit tests do not connect to a database. Integration tests exercise the real database connection with verified TLS, library operations, startup migration/seeding, and HTTP endpoints.

## Run unit tests

Unit tests check isolated configuration, validation, routing, and response-mapping behavior:

```sh
make unit-test
# equivalent:
go test ./... -count=1
```

They need no database credentials or external services. They do not prove that SQL Server connectivity, SQL statements, migrations, or TLS settings work against your instance.

The hardening-report HTTP integration test is separate from the SQL Server suite and does not connect to a database. It starts a local GOREST process to verify the endpoint's default-off behavior, peer-IP allowlisting, forwarded-header handling, report redaction, response headers, and format validation.

## Prepare a SQL Server integration database

Integration tests use `GOREST_DB_*` connection settings from the process environment or `.env`. Before running them:

1. Use a dedicated, non-production SQL Server database and a restricted login. The HTTP flow runs the configured migration and seeder, so the login needs permissions for those operations and the tested user operations.
2. Set `GOREST_DB_DRIVER=sqlserver`, `GOREST_DB_DISABLE=false`, and the host, port, database name, username, and password. The default port is `1433`.
3. Configure `GOREST_DB_TLS=mandatory` or `strict`. Keep `GOREST_DB_TRUST_SERVER_CERTIFICATE=false`; the server certificate must chain to a root trusted by the host and its name must match `GOREST_DB_HOST`. For a private CA, set `GOREST_DB_CA_PATH` to its PEM trust bundle.
4. Set `GOREST_TEST_SQLSERVER_CONFIRM_DATABASE` to exactly match `GOREST_DB_NAME`. This explicit confirmation is required even if the name looks like a test database. The test runner refuses SQL Server system databases and production mode.

The SQL Server integration tests perform schema changes and write test data. Do not point them at production or any database containing data that must be preserved.

## Run integration tests

Run SQL Server TLS, database library, migration/seeder, and HTTP integration tests:

```sh
make test-integration
```

Run the local hardening HTTP integration test independently:

```sh
make test-hardening-integration
```

The SQL Server TLS test connects using mandatory or strict encryption and server-certificate validation. The library test checks user creation, password hashing, listing, and rejection of an incorrect password. The HTTP flow starts the application with migrations and seed data enabled, then checks user creation, duplicate rejection, login, and authenticated listing. Test-created users are uniquely identified and removed afterward; the suite does not truncate tables or drop the database.

To run only the SQL Server TLS connection check:

```sh
make test-sqlserver-tls
```

The `integration` build tag keeps external-service checks out of the default unit-test command. The target above selects SQL Server checks only; it does not run optional Redis tests.

## Optional Redis checks

Redis is not required by the user API and is disabled by default. If you have enabled Redis for your application, the full tagged suite can also test it when you run:

```sh
make integration-reports
```

That command runs the tagged integration suite once and creates coverage and test-output files under `coverage/`. Redis tests use the configured Redis endpoint when enabled, so point it only at a dedicated non-production instance. They use a unique key and do not flush the database. Redis mTLS checks require an endpoint configured to require client certificates. See [Redis configuration](./configuration.md#redis-optional) before enabling Redis.

## Automated checks in GitHub Actions

The workflow runs unit tests and coverage on pushes and pull requests. Its SQL Server integration job runs only when secrets are configured for a reachable, dedicated non-production SQL Server database:

- `GOREST_TEST_SQLSERVER_HOST`, optional `GOREST_TEST_SQLSERVER_PORT` (defaults to `1433`), `GOREST_TEST_SQLSERVER_DATABASE`, `GOREST_TEST_SQLSERVER_USER`, and `GOREST_TEST_SQLSERVER_PASSWORD`.
- `GOREST_TEST_SQLSERVER_TLS`, set to `mandatory` or `strict`.
- `GOREST_TEST_SQLSERVER_CONFIRM_DATABASE`, exactly matching the database name.

The server certificate must validate against the GitHub runner's trusted roots and its hostname must match the configured host. The workflow does not install a private CA bundle; use a certificate trusted by the runner for CI integration tests. If the required secrets are unavailable, GitHub skips the integration job; a skipped job is not a passing test result. Credentials belong in Actions secrets or a secret manager, never in the workflow file.

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
