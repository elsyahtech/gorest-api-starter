# GOREST ScyllaDB API Starter

A ScyllaDB-focused application starter built with [GOREST](https://github.com/elsyahtech/gorest/tree/main). It demonstrates a user API with authentication, password hashing, wide-column persistence, automated unit tests, and opt-in integration tests against a separately confirmed ScyllaDB test keyspace.

## Documentation

Start with the [documentation index](./docs/README.md). The guides cover setup, the active API, architecture, ScyllaDB configuration, and testing.

## Included endpoints

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/users/create` | Create a user and hash its password |
| `POST` | `/users/login` | Verify credentials and issue an access cookie |
| `GET` | `/users/` | List users; requires a valid access cookie |

These are starter examples, not a complete production identity or authorization system. Review the [API reference](./docs/api-reference.md) and secure routes according to your application's requirements.

## Requirements and quick start

- Go `1.26.5` or newer.
- A reachable ScyllaDB server and a pre-created keyspace for the application.

Copy `.env.example` to `.env`, configure the ScyllaDB host, keyspace, optional credentials, and TLS settings, then apply the migration and start the application:

```sh
cp .env.example .env
# Configure the ScyllaDB connection.
GOREST_DB_MIGRATION=true go run .
GOREST_DB_MIGRATION=false go run .
```

The application loads `.env` when present and also supports environment-only startup. ScyllaDB requires the keyspace to exist before connecting. See [Getting started](./docs/getting-started.md) and [ScyllaDB configuration](./docs/database.md) for the full setup.

## Tests and coverage

```sh
make unit-test                  # no ScyllaDB server required
make test-integration           # ScyllaDB connection, migration/seeder, library, and HTTP tests
make test-scylla-connection     # connection-only ScyllaDB integration test
make test-hardening-integration # local HTTP hardening test
make integration-reports        # tagged integration suite and coverage reports
make coverage                   # unit coverage reports
```

Integration tests write to the configured keyspace. Use a dedicated non-production test keyspace ending in `_test` or `_testing`, and set `GOREST_TEST_SCYLLA_CONFIRM_KEYSPACE` to its exact name. For a test keyspace with another name, also set `GOREST_TEST_SCYLLA_ALLOW_NON_TEST_KEYSPACE` to that exact name. Tests apply migrations and seed ten sample rows; they do not drop the keyspace. See the [testing guide](./docs/testing.md). Optional Redis tests may also connect to external Redis services when configured.

GitHub Actions runs unit tests and isolated ScyllaDB integration-configuration checks. It does not connect to ScyllaDB; database integration tests need an explicitly configured test server and keyspace.

## Framework and security

The [GOREST documentation](https://github.com/elsyahtech/gorest/tree/main/docs) is the framework reference. ScyllaDB transactions are not supported by the generic GOREST transaction API; this starter uses single-statement inserts. Enable verified TLS for remote deployments and protect client private keys. Do not commit credentials, keys, or certificates. Redis remains optional and disabled by the starter's `.env.example`; the current user endpoints do not depend on it.

For private vulnerability reporting and the project's security policy, see [SECURITY.md](./SECURITY.md).
