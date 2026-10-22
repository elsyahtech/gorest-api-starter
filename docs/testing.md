# Testing and E2E coverage

The starter's HTTP examples are tested through the registered GOREST routes. These tests exercise the in-process Fiber server, handlers, libraries, models, and a real database adapter; they do not start the application as a separate process.

## Quick commands

Run the default test suite:

```sh
go test ./...
# or
make test
```

Run the E2E package without cached results:

```sh
make test-e2e
```

Run static checks:

```sh
go vet ./...
# or
make check
```

Generate application-code coverage with the default E2E setup:

```sh
make coverage
```

This writes `coverage.out` and `coverage.html`. Coverage artifacts are local build outputs and are ignored by Git. MongoDB and ScyllaDB tests are skipped unless explicitly enabled, so this default report covers the local SQLite E2E paths, not every native driver.

## What the suite covers

The E2E inventory is cross-checked against routes registered in `app/router/`:

- **Departments and orders:** CRUD flows, partial updates, single-resource reads, order lookup by user, plus the registered `HEAD` and `OPTIONS` metadata routes.
- **User ORM queries:** list/detail, selected/all columns, ordering, grouping, distinct, `HAVING`, joins, preload, limit, and offset pagination.
- **User writes:** create variants, upsert, returning results, single and bulk updates, PATCH status, and delete variants. Bulk-update tests include missing and duplicate primary-key validation.
- **Native SQL:** user read/create/update/delete routes are exercised against the E2E SQLite database.
- **Native MongoDB and ScyllaDB:** native read/create/update/delete flows are opt-in and require their respective external services.

The suite verifies representative success and error behavior, response DTOs, password hashing where relevant, and persisted database state. It is not exhaustive proof of every GOREST feature or every driver/version combination. There are currently no separate handler/library/model unit-test suites; the functional examples are primarily verified end to end.

## SQLite E2E setup

The standard test helpers create isolated temporary directories and SQLite databases, copy the relevant migration files, register the same route functions used by the application, and remove the temporary files with Go's test cleanup. They do not use the PostgreSQL database configured for `go run .`. Running `go test ./...` does not require the development PostgreSQL server.

The helpers temporarily change the process working directory so migrations and SQLite files remain isolated. Do not add `t.Parallel()` to tests using these helpers unless the working-directory setup is redesigned to be concurrency-safe.

## Optional native-driver tests

GOREST wires the database integration; neither GOREST nor these tests start Docker or provision database services. Start a disposable database yourself, or use another isolated test service, before enabling the corresponding tests. Do not point destructive E2E tests at a development or production database: tests create and drop a MongoDB database or ScyllaDB keyspace.

### MongoDB

The test helper defaults to `127.0.0.1:27018`. It currently connects without MongoDB credentials, so use a local, non-production test instance with authentication disabled and bind it to loopback. Override the host or port when needed:

```sh
GOREST_MONGO_HOST=127.0.0.1 GOREST_MONGO_PORT=27018 make test-mongo
```

`make test-mongo` sets `GOREST_MONGO_E2E=1`. The test account/service must permit creating and dropping the test databases.

### ScyllaDB

The test helper defaults to `127.0.0.1:9042`. It creates a temporary keyspace and `users` table, then drops the keyspace during cleanup. The configured account must have permission to create and drop keyspaces. Optional credentials are read from `GOREST_SCYLLA_USERNAME` and `GOREST_SCYLLA_PASSWORD`:

```sh
GOREST_SCYLLA_HOST=127.0.0.1 \
GOREST_SCYLLA_PORT=9042 \
GOREST_SCYLLA_USERNAME=test-user \
GOREST_SCYLLA_PASSWORD=test-password \
make test-scylla
```

`make test-scylla` sets `GOREST_SCYLLA_E2E=1`.

Run the full E2E suite with both native drivers enabled and refresh the combined coverage report:

```sh
GOREST_MONGO_HOST=127.0.0.1 GOREST_MONGO_PORT=27018 \
GOREST_SCYLLA_HOST=127.0.0.1 GOREST_SCYLLA_PORT=9042 \
GOREST_SCYLLA_USERNAME=test-user GOREST_SCYLLA_PASSWORD=test-password \
make coverage-e2e
```

This target sets both opt-in flags. Both services must be ready and reachable; MongoDB's current test helper has no credentials option. The command runs all E2E tests and writes `coverage.out` and `coverage.html`, including the native-driver code paths exercised by those tests.

## Important scope limits

- Native SQL E2E uses SQLite. Passing those tests confirms the starter's SQLite test path; it does not prove that PostgreSQL, MySQL, Oracle, or SQL Server syntax and behavior are equivalent. The native SQL examples use PostgreSQL-style placeholders and must be adapted and tested against the intended driver.
- ORM behavior, schema features, constraints, and transaction support vary between SQL databases, MongoDB, and ScyllaDB. Transaction calls through GOREST do not provide SQL-style transaction atomicity for MongoDB or ScyllaDB.
- MongoDB and ScyllaDB E2E test only the service and driver versions you run locally. They do not certify a production topology, replica set, authentication policy, TLS setup, or cluster behavior.
- Coverage percentage measures which application statements the tests executed; it is not a measure of test quality, production readiness, or cross-driver compatibility.
- The application has no authentication/authorization on `/examples`. Protect and adapt routes before exposing them to clients.

See the [API reference](./api-reference.md) for route and payload behavior, the [database guide](./database.md) for driver and transaction notes, and the [configuration guide](./configuration.md) before deployment.
