# Testing your Oracle starter

The starter provides fast unit tests and opt-in Oracle integration tests. Unit tests do not connect to a database. Integration tests exercise an Oracle connection, migrations, seeders, library operations, and HTTP endpoints.

## Run unit tests

```sh
make unit-test
# equivalent:
go test ./... -count=1
```

Unit tests need no database credentials or external services. They do not prove Oracle connectivity, SQL syntax, migrations, or wallet configuration.

The hardening-report HTTP integration test is separate from the Oracle suite and does not connect to a database. It starts a local GOREST process to verify the endpoint's default-off behavior, peer-IP allowlisting, forwarded-header handling, report redaction, response headers, and format validation.

## Prepare an Oracle integration schema

Oracle tests require a compatible Oracle Instant Client and a dedicated, disposable, non-production application schema. They perform schema changes and seed data; do not point them at production or a schema containing data that must be preserved.

The schema must be able to create the GOREST migration lock and migration/seeder history tables, the starter `users` table, and its index. Have the database administrator grant `CREATE TABLE` and a quota on the schema's target tablespace. Without these, Oracle connectivity may pass while migration and HTTP tests fail.

1. Configure `GOREST_DB_DRIVER=oracle`, `GOREST_DB_DISABLE=false`, `GOREST_DB_HOST`, `GOREST_DB_PORT` (default `1521`), `GOREST_DB_NAME` (the service name, not SID), `GOREST_DB_USER`, and `GOREST_DB_PASSWORD`.
2. Set `GOREST_DB_LIB_DIR_PATH` if Oracle Instant Client is not in the system library path.
3. For TCPS, set `GOREST_DB_TLS=enable` and `GOREST_DB_CA_PATH` to a valid Oracle wallet directory. `GOREST_DB_CA_PATH` is a wallet/TNS_ADMIN directory, not a PEM file. For isolated test servers without TCPS, `disable` is accepted; do not use unencrypted TCP over untrusted networks.
4. Set `GOREST_TEST_ORACLE_CONFIRM_USER` to exactly match the dedicated schema username in `GOREST_DB_USER`. The test runner refuses to run without this exact confirmation, in production mode, or as a built-in administrative account.

The tests do not truncate tables or drop the schema. They do apply migrations and leave the migration/seeder tracking objects and ten sample users in the confirmed schema. The HTTP test writes a unique user and removes it afterward; the library test similarly removes its unique user.

## Run Oracle integration tests

```sh
make test-integration
```

This runs the Oracle config-safety tests, a connection check, library persistence tests, and the HTTP flow. To run only the Oracle connection check:

```sh
make test-oracle-connection
```

Run the local hardening HTTP integration test independently:

```sh
make test-hardening-integration
```

The `integration` build tag keeps external-service checks out of the default unit-test command.

## Optional Redis checks

Redis is not required by the user API and is disabled by default. `make integration-reports` runs the full tagged suite and creates coverage and test-output files under `coverage/`. The Oracle tests use the confirmed external schema, and optional Redis checks may connect to the configured endpoint if Redis is enabled. Point it only at a dedicated non-production instance. Redis checks use unique keys and do not flush the database; Redis mTLS checks require a server configured to require client certificates. See [Redis configuration](./configuration.md#redis-optional) before enabling Redis.

## Automated checks in GitHub Actions

GitHub Actions runs unit tests and the isolated Oracle integration-configuration tests. It does not connect to Oracle: full integration tests require a compatible Oracle Instant Client, wallet when using TCPS, and a dedicated test schema.

## Coverage

Generate unit coverage locally with:

```sh
make coverage
```

Generate integration test output and application coverage with:

```sh
make integration-reports
```

Reports are written under `coverage/`. They help show which code paths tests exercised; they do not replace reviewing the behavior and test cases relevant to your application.
