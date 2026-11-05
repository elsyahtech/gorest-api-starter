# GOREST MySQL API Starter

A MySQL-focused application starter built with [GOREST](https://github.com/elsyahtech/gorest/tree/main). It demonstrates a small user API with authentication, password hashing, a MySQL data-access layer, automated unit tests, and opt-in HTTP integration tests against a separately configured MySQL test database.

## Documentation

Start with the [documentation index](./docs/README.md). The guides cover setup, the active API, architecture, MySQL configuration, and testing evidence.

## Included endpoints

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/users/create` | Create a user and hash its password |
| `POST` | `/users/login` | Verify credentials and issue an access cookie |
| `GET` | `/users/` | List users; requires a valid access cookie |

These are starter examples, not a complete production identity or authorization system. Review the [API reference](./docs/api-reference.md) and secure routes according to your application's requirements.

## Requirements and quick start

- Go `1.26.5` or newer.
- A reachable MySQL server and a database for the application.

Copy `.env.example` to `.env`, provide the MySQL database name, username, password, and appropriate TLS settings, then apply the MySQL migration and start the application:

```sh
cp .env.example .env
# Edit .env with your MySQL settings.
GOREST_DB_MIGRATION=true go run .
GOREST_DB_MIGRATION=false go run .
```

The application loads `.env` when present and also supports environment-only startup. See [Getting started](./docs/getting-started.md) and [MySQL configuration](./docs/database.md) for the full setup.

## Tests and coverage

```sh
make unit-test       # isolated unit tests; no database required
make coverage        # unit coverage profile, text summary, and HTML report
make test-integration # E2E requests against the explicitly configured MySQL test host
make integration-reports # one E2E run: readable and JSON reports plus app coverage
```

Integration tests are opt-in and use the test-specific `GOREST_TEST_MYSQL_*` settings in the same `.env` file (or environment variables). They require a dedicated test database, normally ending in `_test` or `_testing`, and exact database-name confirmation before making test-data changes. The optional Redis integration test uses the `GOREST_REDIS_*` settings only when Redis is explicitly enabled; point these at a non-production test instance. Tests do not use Docker or a disposable local database. Read the [testing guide](./docs/testing.md) before running them.

The coverage outputs and integration reports are generated locally under `coverage/`; rerun the documented commands to regenerate evidence for the current source and environment. The integration report command runs E2E once and generates both report formats and an application coverage profile.

GitHub Actions runs unit tests and coverage on pushes and pull requests. Its MySQL E2E job runs when the required external test database secrets are configured; see the [testing guide](./docs/testing.md#github-actions-ci).

## Framework and security

The [GOREST documentation](https://github.com/elsyahtech/gorest/tree/main/docs) is the framework reference. For MySQL TLS, private-CA, and mTLS setup, follow the [MySQL configuration guide](./docs/database.md#configure-tls-private-cas-and-client-certificates) and verify that the chosen connector supports your database's certificate policy. Do not commit credentials, private keys, or certificates. The default Redis wiring remains optional and disabled by the starter's `.env.example`; the current user endpoints do not depend on Redis.

For private vulnerability reporting and the project's security policy, see [SECURITY.md](./SECURITY.md).
