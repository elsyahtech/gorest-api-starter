# Testing your MongoDB starter

The starter provides isolated unit tests and opt-in MongoDB integration tests. Unit tests do not connect to a database. Integration tests exercise MongoDB connectivity, migrations, indexes, user-library operations, and HTTP endpoints. Redis checks are separate optional tests.

## Run unit tests

```sh
make unit-test
# equivalent:
go test ./... -count=1
```

Unit tests need no database credentials or external services. They do not prove MongoDB connectivity, document validation, index creation, TLS, or server topology behavior.

The hardening-report HTTP integration test is separate from MongoDB and runs locally:

```sh
make test-hardening-integration
```

## Prepare a MongoDB integration database

Integration tests run migrations and write user documents. Use a dedicated, disposable, non-production MongoDB database; never use production or a database containing data that must be preserved.

1. Prefer explicit `GOREST_TEST_MONGO_HOST`, `PORT`, `DATABASE`, `USER`, `PASSWORD`, `AUTH_SOURCE`, and `TLS` values. If all dedicated connection fields are blank, the test may fall back to `GOREST_DB_*` only when `GOREST_APP_ENV` is `development`, `test`, or `testing`, the driver is `mongo`, and the database is enabled.
2. Set `GOREST_TEST_MONGO_CONFIRM_DATABASE` to exactly match the target database. This is required before any integration test connects.
3. The name must end in `_test` or `_testing`. If the dedicated non-production database has another name, also set `GOREST_TEST_MONGO_ALLOW_NON_TEST_DATABASE` to that exact name.
4. Configure `GOREST_TEST_MONGO_TLS=true` for a TLS endpoint, with CA and optional client-certificate/key files as needed. Client certificate and key must be supplied together.
5. Grant the test account only the permissions needed to create/access the GOREST migration metadata, `users` collection, indexes, and documents.

The suite creates the collection and indexes if needed, records migration state, and runs the sample seeder once. Three sample documents with placeholder password hashes and migration/seeder history remain after testing. Test cases insert UUID-scoped users and delete only their own IDs/emails. The suite does not drop collections; do not run the migration down file because it drops the whole `users` collection.

The starter constructs a `mongodb://host:port` URI; MongoDB SRV URI configuration is not available in these settings. `GOREST_DB_AUTH_SOURCE` defaults to `admin`; credentials are optional but must be configured as a pair.

## Run MongoDB integration tests

```sh
make test-integration
```

This runs MongoDB config-safety tests, a connection check, library persistence tests, and an HTTP flow. For the connection-only test:

```sh
make test-mongo-connection
```

The `integration` build tag keeps external database tests out of the default unit-test command. The safety checks fail before connecting if confirmation or connection settings are missing.

## Optional Redis checks and reports

Redis is not required by the user API and is disabled by default. `make integration-reports` runs the full tagged suite and writes test and coverage reports under `coverage/`. MongoDB integration checks use the confirmed test database. Redis checks may connect to an external Redis instance if it is enabled in the environment; configure a dedicated non-production endpoint only. Redis tests use unique keys and do not flush the database. Redis mTLS checks require a server configured to require client certificates.

`GOREST_REDIS_DISABLE=true` skips the regular Redis integration check, but does not disable the separate Redis mTLS test when its dedicated test settings or runtime client certificate/key are configured. To keep the report run focused on MongoDB and local HTTP checks, also leave `GOREST_TEST_REDIS_MTLS_*` and the runtime Redis client certificate/key paths unset:

```sh
GOREST_REDIS_DISABLE=true make integration-reports
```

## GitHub Actions

GitHub Actions runs unit tests and isolated MongoDB integration-configuration checks. It does not connect to MongoDB. Full integration tests require an explicitly configured server, TLS material if needed, a dedicated test database, and exact database confirmation.

## Coverage

Generate unit coverage locally:

```sh
make coverage
```

Generate integration output and application coverage:

```sh
make integration-reports
```

Reports are written under `coverage/`. They help show which code paths tests exercised; they do not replace reviewing the behavior and test cases relevant to your application.
