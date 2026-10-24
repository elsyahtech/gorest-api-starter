# Testing and E2E coverage

The starter's HTTP examples are tested through the registered GOREST routes. These tests exercise the in-process Fiber server, handlers, libraries, models, and a real database adapter; they do not start the application as a separate process.

## Quick commands

Run the default test suite:

```sh
go test ./...
# or
make test
```

`make test` runs `go test -v ./... -count=1`, so it prints each test as it runs, disables cached results, and includes all packages and lifecycle tests even when their names do not contain `E2E`. A final success message appears only after the full command passes.

This is the regular Go suite and does not launch external database services. Driver tests in `./test/drivers` skip unless a driver is explicitly selected. To run the isolated seven-driver matrix, use `make test-e2e`; it requires Docker and more time/resources.

Run the complete database verification matrix against disposable database containers:

```sh
make test-e2e
```

`make test-e2e` runs the in-process route suite on temporary SQLite, then runs driver integration and route scenarios for all seven database drivers sequentially. Each target first runs the shared ORM route/lifecycle scenarios that the driver supports, then runs that driver's native query routes in a separate stage. PostgreSQL, MySQL, Oracle, and SQL Server run the complete shared suite, including native SQL routes. SQLite runs the shared suite against a temporary database. MongoDB runs the applicable shared ORM routes and lifecycle scenarios, skips ORM transaction-atomicity checks (not provided by the exercised GOREST API), and treats its unsupported ORM `Join()` route as a passing HTTP 400 validation scenario before running the four native MongoDB routes. ScyllaDB runs the shared ORM routes supported by its CQL schema/API, reports the unsupported SQL-style operations it omits, and then runs its four native CQL routes. The runner prints these exclusions before the shared suite so the output explains why driver logs differ. The SQLite-specific target filters out MongoDB and ScyllaDB route tests, so they do not appear as unrelated skips in its output. The runner creates isolated services and removes them after each driver. To run only one part of the same matrix, use `make test-e2e-sqlite`, `make test-e2e-postgres`, `make test-e2e-mysql`, `make test-e2e-oracle`, `make test-e2e-sqlserver`, `make test-e2e-mongo`, or `make test-e2e-scylla`. The runner requires Docker Compose v2. It does not use the application's database host, database name, or credentials.

Redis has a separate disposable E2E target because it is an optional cache service, not one of the seven database drivers:

```sh
make test-e2e-redis
```

This target starts the Redis service from `test/drivers/compose.yaml`, tests Redis user CRUD and conflict handling, and removes its container and network afterward. It does not use the application's Redis service or configured host.

Run the pinned baseline linters or all configured linters:

```sh
make lint       # errcheck, govet, staticcheck, ineffassign, unused
make lint-full  # every configured golangci-lint check
```

`make lint-full` may report existing quality-backlog findings. The release gate uses the pinned baseline `make lint` target; see the [release checks guide](./release-checks.md).

Run static checks:

```sh
go vet ./...
# or
make check
```

`make check` runs `make vet` followed by `make test`, and prints a status message for each stage and a final success message.

For the pinned release gate, including lint, secret scanning, and the SQLite driver integration test, see the [release checks guide](./release-checks.md).

## Production configuration scenarios

The default Go test suite also checks production configuration rules without connecting to production services. These tests set `APP_ENV=production` in an isolated test process and call the configuration loaders directly. They verify required HTTP bind and CORS settings, required server-database settings, rejection of automatic migrations and seeders, the explicit directory requirement for production SQLite, and the Redis host/password/TLS requirements when Redis is enabled. They also check that valid production settings are accepted.

Run them with the ordinary suite:

```sh
go test ./...
```

These are configuration validation tests; they do not start the server or connect to a production database or Redis service.

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
- **Native SQL:** user read/create/update/delete routes run against all five SQL drivers through the seven-driver runner; the ordinary `make test` suite exercises SQLite.
- **Native MongoDB and ScyllaDB:** native read/create/update/delete routes run against disposable MongoDB and ScyllaDB containers in `make test-e2e`; the ordinary `make test` suite skips these external-service tests.

The suite verifies representative success and error behavior, response DTOs, password hashing where relevant, and persisted database state. It is not exhaustive proof of every GOREST feature or every driver/version combination. There are currently no separate handler/library/model unit-test suites; the functional examples are primarily verified end to end.

## Driver verification matrix

The automated driver integration tests run startup migrations and seeders, ORM CRUD, and transaction rollback where supported. Evidence pages record earlier runs as well as supplemental runs of the current per-driver targets. Oracle's repeat runs include a test-harness timeout adjustment recorded on its evidence page. The manual column was not generated by `go test`. Select a driver's result link for its recorded test output and target details.

`passed` means the linked run met the assertion recorded for that scenario. `partial` means only the explicitly described portion was exercised. `not tested` means there is no run evidence; `not applicable` means the feature does not apply to that driver/API. MongoDB and ScyllaDB transaction atomicity is not applicable to the GOREST ORM transaction API exercised here. The evidence is scoped to the exact server versions and configurations listed on each result page.

| Driver | Manual actual verification | Migration E2E | Seeder E2E | CRUD E2E | Transaction E2E | Mid-step rollback | Conflicts / invalid operations | Concurrency | Startup failure | Shutdown | Shared ORM route/lifecycle E2E | Native route E2E (2026-10-24) | Evidence |
| :----- | :------------------------- | :------------ | :--------- | :------- | :-------------- | :---------------- | :---------------------------- | :--------- | :-------------- | :------- | :-------------------------- | :---------------------------- | :------- |
| PostgreSQL | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` (SQL) | [Test Result](./testing-result/testing-result-postgres.md) |
| MySQL | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` (SQL) | [Test Result](./testing-result/testing-result-mysql.md) |
| Oracle | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` (SQL) | [Test Result](./testing-result/testing-result-oracle.md) |
| SQL Server | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` (SQL) | [Test Result](./testing-result/testing-result-sqlserver.md) |
| SQLite | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` | `passed` (SQL) | [Test Result](./testing-result/testing-result-sqlite.md) |
| MongoDB | `passed` | `passed` | `passed` | `passed` | `not applicable` | `not applicable` | `passed` | `passed` | `passed` | `passed` | `passed` (applicable scenarios; transactions N/A) | `passed` (MongoDB) | [Test Result](./testing-result/testing-result-mongo.md) |
| ScyllaDB | `passed` | `passed` | `passed` | `passed` | `not applicable` | `not applicable` | `passed` | `passed` | `passed` | `passed` | `partial` (CQL-supported scenarios; exclusions printed by runner) | `passed` (CQL) | [Test Result](./testing-result/testing-result-scylla.md) |

The `test/drivers` integration test exercises startup migrations and seeders, ORM CRUD, native SQL route CRUD for SQL drivers, duplicate primary keys, rejected requests, concurrent HTTP creates, invalid-migration startup failure, and normal `app.Close()` cleanup against each configured test driver. For SQL drivers, it also verifies that a database constraint failure on the second transaction write rolls back the first write. MongoDB and ScyllaDB do not use the ORM transaction API for atomic rollback, so transaction and mid-step rollback are not applicable. Native MongoDB and CQL route operations run from `test/e2e` while their disposable services are active. See the [driver test result index](./testing-result.md) for previously recorded evidence; rerun `make test-e2e` to generate current results for this working tree.

Concurrency coverage is a bounded batch of eight simultaneous valid creates. The assertion requires all eight responses to be HTTP `200`, contain `success=true`, and match persisted rows; any transport error, 4xx, or 5xx fails the run. All seven drivers passed this assertion. Oracle passed three consecutive reruns after the Fiber test-client timeout was aligned with the configured database operation timeout. The SQLite run uses a single-connection database pool, which queues writes for SQLite's single-writer model.

Lifecycle evidence covers both startup failure paths and graceful shutdown. For PostgreSQL, MySQL, Oracle, SQL Server, MongoDB, and ScyllaDB, a subprocess test points the driver at a reserved then released local port and asserts that GOREST exits with `failed initializing database`; each driver integration run also rejects an invalid migration after connecting. SQLite is embedded, so an unavailable network endpoint does not apply; its invalid-migration startup failure is tested directly. Each driver integration run verifies successful `app.Close()` for its database and logger, and the app-level E2E test verifies that shutdown waits for an in-flight HTTP request before exiting. The shutdown contract tested here is graceful completion within GOREST's configured shutdown period; forced-timeout behavior is outside this scenario.

Use disposable, empty databases or schemas for DB-backed tests. Startup seeders insert fixed sample IDs, so repeating the integration test against a populated target may conflict. SQLite uses a temporary file. Oracle integration requires a test-only service, an available Instant Client, and explicit connection settings through the `GOREST_DRIVER_*` environment variables; the SQLite default in `app/config/database.go` does not configure Oracle credentials. The recorded driver run used existing local test services; its command and target details are in [testing-result.md](./testing-result.md).

## Seven-driver E2E matrix

`test/drivers/database_integration_test.go` contains the shared integration scenarios: startup migrations and seeders, ORM CRUD, transaction rollback for SQL drivers, conflicts and invalid operations, concurrent HTTP requests, connection and migration startup failures, and application shutdown. The native SQL route scenarios run against SQLite, PostgreSQL, MySQL, Oracle, and SQL Server. Native MongoDB route scenarios run against MongoDB; native CQL route scenarios run against ScyllaDB. These scenarios use the in-process HTTP server and real database services in disposable containers.

On October 24, 2026, the targeted `make test-e2e-oracle`, `make test-e2e-sqlserver`, `make test-e2e-mongo`, and `make test-e2e-scylla` runs all passed. Oracle and SQL Server passed the full shared route/lifecycle suite. MongoDB passed its applicable shared ORM route/lifecycle scenarios and all four native MongoDB route scenarios. ScyllaDB passed the shared scenarios supported by its CQL schema/API and all four native CQL route scenarios; unsupported SQL-style operations and transaction checks were explicitly omitted. PostgreSQL and MySQL full shared route/lifecycle runs also passed on this date. Each target used an isolated disposable container and removed its container and network afterward. Earlier Oracle container startup attempts failed under memory pressure, but the current Oracle target completed successfully. See each driver's result page for details.

The database integration testing documented here has been run on macOS and Linux. Windows is not currently supported or included in the verified testing scope for the automation; the runner uses Bash and OS-specific Oracle Instant Client packages.

### Run all seven drivers or one driver

Requirements: Go, Bash, Docker Engine with Docker Compose v2, and network access the first time the database images and Oracle Instant Client are downloaded. Linux also needs `curl`, `unzip`, and `sha256sum` or `shasum`; macOS needs `curl`, `hdiutil`, `ditto`, and `shasum`. Windows native and WSL runs are outside the documented support scope.

```sh
make test-e2e
```

The runner executes these steps sequentially:

1. Run the full in-process route suite and database integration checks on temporary SQLite.
2. Start PostgreSQL, MySQL, SQL Server, MongoDB, and ScyllaDB one at a time, run that driver's matrix scenarios, and remove its container.
3. Start Oracle Free, run Oracle matrix scenarios using the `gorest_test` schema in `FREEPDB1`, and remove its container.
4. Remove the Compose project network, volumes, and temporary Oracle Instant Client files, including after test failure. On failure, the runner prints logs for the service that was running.

Choose a single driver with its corresponding target, for example:

```sh
make test-e2e-sqlserver
```

The available suffixes are `sqlite`, `postgres`, `mysql`, `oracle`, `sqlserver`, `mongo`, and `scylla`.

Each container receives an available loopback port, and the Compose project name is unique per run. This avoids fixed-port conflicts and leaves other projects and database containers alone. Database files live only in each disposable container; Compose images remain in Docker's local cache after the run. SQLite's test database and any downloaded Oracle Instant Client are stored temporarily on the host and removed when the script exits. The host connection settings in `GOREST_DRIVER_*` for your normal local databases do not redirect this automated workflow to those services.

Oracle uses the community-maintained `gvenzl/oracle-free:23.26.3-slim` image, which runs Oracle Database Free, creates the app schema using `APP_USER`/`APP_USER_PASSWORD`, and provides a health check. Its published image supports both Intel and ARM Linux containers. Oracle Free has edition limits, so the result verifies this driver against that edition, not every Oracle deployment. See the [image documentation](https://hub.docker.com/r/gvenzl/oracle-free). Oracle needs at least 2 GB of memory and can take several minutes to initialize; Docker Desktop must have enough memory assigned. SQL Server and ScyllaDB are configured for `linux/amd64`, so ARM hosts rely on Docker's architecture emulation for those two services. ScyllaDB runs with one CPU, 1 GB memory, and developer mode; SQL Server uses the Developer edition.

Oracle's Go driver also needs Oracle Instant Client on the host running the tests. If `GOREST_DRIVER_LIB_DIR` points to an existing compatible Instant Client directory, the runner uses it without changing or removing it. Otherwise it downloads a pinned Basic package, verifies its SHA-256 checksum, uses it for the test, and removes it at exit:

| Host running tests | Instant Client package |
| :----------------- | :--------------------- |
| macOS Intel | 19.16 DMG |
| macOS ARM64 | 23.26 DMG |
| Linux x86_64 | 23.26 ZIP |
| Linux ARM64 | 23.26 ZIP |

The Oracle image, Oracle Instant Client, and application integration test are separate pieces: the image supplies the database service; the client library is loaded by the Go test process on the host. The runner does not require a local Oracle service or Oracle credentials supplied by the user.

Optional settings for the automated target:

| Variable | Default | Purpose |
| :------- | :------ | :------ |
| `GOREST_TEST_DRIVERS_PASSWORD` | `GorestTest!Pass2026` | Password for test-only users in the disposable database containers. Set a different value if local policy requires it. |
| `GOREST_TEST_DRIVERS_ORACLE_USERNAME` | `gorest_test` | Oracle schema user created in the temporary Oracle Free service. It uses the password above. |
| `GOREST_DRIVER_LIB_DIR` | downloaded into temporary storage | Use an existing host Instant Client directory. The runner does not remove or modify this directory. |
| `GOREST_TEST_DRIVERS_SKIP_ORACLE` | `0` | Set to `1` to omit Oracle. The result then covers SQLite plus five server drivers, not all seven drivers. |

To run SQLite and the five non-Oracle containerized drivers without Oracle:

```sh
GOREST_TEST_DRIVERS_SKIP_ORACLE=1 make test-e2e
```

For lower-level development, the Go integration test can still be run directly against a database service you provide. This is an advanced path; unlike `make test-e2e`, it does not provision or isolate that service. Use a disposable database and configure `GOREST_DRIVER_INTEGRATION` and connection settings explicitly:

```sh
GOREST_DRIVER_INTEGRATION=postgres \
go test ./test/drivers -run '^TestDatabaseDriverIntegration$' -count=1 -v
```

Supported names and defaults are:

| Driver | Value | Default host and port | Default database/service |
| :----- | :---- | :-------------------- | :----------------------- |
| PostgreSQL | `postgres` | `127.0.0.1:15432` | `gorest_driver_it` |
| MySQL | `mysql` | `127.0.0.1:13306` | `gorest_driver_it` |
| Oracle | `oracle` | `127.0.0.1:1521` | `gorest_driver_it` |
| SQL Server | `sqlserver` | `127.0.0.1:11433` | `gorest_driver_it` (created from `master` if missing) |
| SQLite | `sqlite` | temporary local file | created per test run |
| MongoDB | `mongo` | `127.0.0.1:27018` | `gorest_driver_it` |
| ScyllaDB | `scylla` | `127.0.0.1:19042` | `gorest_driver_it` keyspace |

For server databases, `GOREST_DRIVER_HOST`, `GOREST_DRIVER_PORT`, `GOREST_DRIVER_DATABASE`, `GOREST_DRIVER_USERNAME`, and `GOREST_DRIVER_PASSWORD` override the defaults. For example, run PostgreSQL with a local test account like this:

```sh
GOREST_DRIVER_INTEGRATION=postgres \
GOREST_DRIVER_HOST=127.0.0.1 \
GOREST_DRIVER_PORT=15432 \
GOREST_DRIVER_DATABASE=gorest_driver_it \
GOREST_DRIVER_USERNAME=gorest_test \
GOREST_DRIVER_PASSWORD='<test-only-password>' \
go test ./test/drivers -run '^TestDatabaseDriverIntegration$' -count=1 -v
```

The same test command works for each driver by changing the integration value. MySQL, Oracle, SQL Server, MongoDB, and ScyllaDB can use the shared connection overrides above. Oracle also reads `GOREST_DRIVER_LIB_DIR` for the Oracle Instant Client library directory. MongoDB and ScyllaDB credentials are optional in the helper, but the test service must permit the operations required by the test. SQL Server needs permission to create the test database when it does not already exist. SQLite needs no external service and ignores the server connection settings.

The direct test command does not start or stop containers. It runs migrations and seeders and writes test data, so use disposable databases or schemas. For the supported automated workflow, use `make test-e2e` or a driver-specific target instead.

To run only the startup-failure check without connecting to a working database, set `GOREST_DRIVER_LIFECYCLE_DRIVER` to `postgres`, `mysql`, `oracle`, `sqlserver`, `mongo`, or `scylla`. The test starts a subprocess against a reserved then released local port and requires GOREST initialization to fail. Oracle also needs `GOREST_DRIVER_LIB_DIR`. This check is not applicable to embedded SQLite; its invalid-migration startup failure is covered by the regular integration test. The full `test/e2e` lifecycle tests separately verify database connection failure, migration failure, and graceful shutdown while a request is active.

The automated target runs the seven-driver matrix sequentially. A successful run proves only the selected image/server versions and configuration used by that run; record those details with the output when publishing evidence.

## SQLite E2E setup

The standard test helpers create isolated temporary directories and SQLite databases, copy the relevant migration files, register the same route functions used by the application, and remove the temporary files with Go's test cleanup. They do not use the persistent SQLite database created by `go run .`. Running `go test ./...` does not require an external database server.

The helpers temporarily change the process working directory so migrations and SQLite files remain isolated. Do not add `t.Parallel()` to tests using these helpers unless the working-directory setup is redesigned to be concurrency-safe.

## Native driver routes

The automated runner enables MongoDB and ScyllaDB route tests while their disposable containers are running. MongoDB tests create and drop uniquely named temporary databases. ScyllaDB tests create and drop uniquely named keyspaces using `NetworkTopologyStrategy` with replication factor 1 for `datacenter1`. No local database service or application credential is used by these targets.

`make coverage` remains the lightweight SQLite route coverage command. `make coverage-e2e` is a manual coverage option for native MongoDB and ScyllaDB paths; it requires both services and their test connection settings to be configured separately.

## Important scope limits

- Native SQL route scenarios run on all five SQL drivers through `make test-e2e`. SQL syntax and behavior can still differ by server version and configuration.
- ORM behavior, schema features, constraints, and transaction support vary between SQL databases, MongoDB, and ScyllaDB. Transaction calls through GOREST do not provide SQL-style transaction atomicity for MongoDB or ScyllaDB.
- MongoDB and ScyllaDB E2E test only the service and driver versions you run locally. They do not certify a production topology, replica set, authentication policy, TLS setup, or cluster behavior.
- Coverage percentage measures which application statements the tests executed; it is not a measure of test quality, production readiness, or cross-driver compatibility.
- The application has no authentication/authorization on `/examples`. Protect and adapt routes before exposing them to clients.

See the [API reference](./api-reference.md) for route and payload behavior, the [database guide](./database.md) for driver and transaction notes, and the [configuration guide](./configuration.md) before deployment.
