# Testing your SQLite starter

The starter provides fast unit tests and SQLite integration tests. Unit tests do not connect to a database. Integration tests exercise actual SQLite persistence, migrations, seeders, library operations, and HTTP endpoints using temporary database files.

## Run unit tests

```sh
make unit-test
# equivalent:
go test ./... -count=1
```

Unit tests need no database credentials or external services. They do not prove the integration behavior covered by the tagged tests.

The hardening-report HTTP integration test is separate from the SQLite suite and does not connect to a database. It starts a local GOREST process to verify the endpoint's default-off behavior, peer-IP allowlisting, forwarded-header handling, report redaction, response headers, and format validation.

## Run SQLite integration tests

```sh
make test-integration
```

Each test creates an isolated SQLite database in a temporary directory. Migrations and seeders come from `app/database/migrations/sqlite` and `app/database/seeders/sqlite`; the normal `app/database/gorest.db` is not used. No database server, credentials, Docker, or `.env` database settings are needed.

The suite checks that the SQLite test configuration resolves to the repository migration/seeder assets, verifies user creation and password hashing through the library, and tests HTTP creation, duplicate rejection, login, and authenticated listing. Temporary files are cleaned up by Go's testing framework.

Run the local hardening HTTP integration test independently:

```sh
make test-hardening-integration
```

The `integration` build tag keeps integration checks out of the default unit-test command.

## Optional Redis checks

Redis is not required by the user API and is disabled by default. `make integration-reports` runs the full tagged suite and creates coverage and test-output files under `coverage/`. The Redis integration checks are skipped unless Redis is explicitly enabled/configured; if enabled, they connect to the configured endpoint, so use a dedicated non-production instance. They use unique test keys and do not flush the database. Redis mTLS checks require an endpoint configured to require client certificates. See [Redis configuration](./configuration.md#redis-optional) before enabling Redis.

## Automated checks in GitHub Actions

GitHub Actions runs unit tests and SQLite integration tests on pushes and pull requests. The integration workflow uses temporary local database files and requires no external database secrets. Optional Redis checks remain disabled unless explicitly configured.

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
