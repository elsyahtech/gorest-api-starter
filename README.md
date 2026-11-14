# GOREST MongoDB API Starter

A MongoDB-focused application starter built with [GOREST](https://github.com/elsyahtech/gorest/tree/main). It demonstrates a user API with authentication, password hashing, document persistence, automated unit tests, and opt-in integration tests against a separately confirmed MongoDB test database.

## Documentation

Start with the [documentation index](./docs/README.md). The guides cover setup, the active API, architecture, MongoDB configuration, and testing.

## Included endpoints

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/users/create` | Create a user and hash its password |
| `POST` | `/users/login` | Verify credentials and issue an access cookie |
| `GET` | `/users/` | List users; requires a valid access cookie |

These are starter examples, not a complete production identity or authorization system. Review the [API reference](./docs/api-reference.md) and secure routes according to your application's requirements.

## Requirements and quick start

- Go `1.26.5` or newer.
- A reachable MongoDB server and a database for the application.

Copy `.env.example` to `.env`, configure the MongoDB host, database name, optional credentials/authentication database, and TLS settings, then apply the migration and start the application:

```sh
cp .env.example .env
# Configure the MongoDB connection.
GOREST_DB_MIGRATION=true go run .
GOREST_DB_MIGRATION=false go run .
```

This starter uses a `mongodb://host:port` connection, not an SRV URI. `GOREST_DB_AUTH_SOURCE` defaults to `admin`; MongoDB username and password must be configured together. For remote deployments, enable TLS and verify the server certificate. See [Getting started](./docs/getting-started.md) and [MongoDB configuration](./docs/database.md).

## Tests and coverage

```sh
make unit-test                  # no MongoDB server required
make test-integration           # MongoDB connection, migration/seeder, library, and HTTP tests
make test-mongo-connection      # connection-only MongoDB integration test
make test-hardening-integration # local HTTP hardening test
make integration-reports        # tagged integration suite and coverage reports
make coverage                   # unit coverage reports
```

MongoDB integration tests write to the configured database. Use a dedicated non-production test database ending in `_test` or `_testing`, and set `GOREST_TEST_MONGO_CONFIRM_DATABASE` to its exact name. For a test database with another name, also set `GOREST_TEST_MONGO_ALLOW_NON_TEST_DATABASE` to that exact name. See the [testing guide](./docs/testing.md). The regular Redis check runs only when enabled; the separate Redis mTLS test also runs when its test settings or runtime client certificate/key are configured, regardless of `GOREST_REDIS_DISABLE`.

GitHub Actions runs unit tests and isolated MongoDB integration-configuration checks. It does not connect to MongoDB; database integration tests need an explicitly configured test server.

## Framework and security

The [GOREST documentation](https://github.com/elsyahtech/gorest/tree/main/docs) is the framework reference. Do not commit credentials, private keys, or certificates. Redis remains optional and disabled by the starter's `.env.example`; the current user endpoints do not depend on it.

For private vulnerability reporting and the project's security policy, see [SECURITY.md](./SECURITY.md).
