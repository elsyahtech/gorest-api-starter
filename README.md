# GOREST SQLite API Starter

A SQLite-focused application starter built with [GOREST](https://github.com/elsyahtech/gorest/tree/main). It demonstrates a small user API with authentication, password hashing, SQLite persistence, automated unit tests, and isolated integration tests.

## Documentation

Start with the [documentation index](./docs/README.md). The guides cover setup, the active API, architecture, SQLite configuration, and testing.

## Included endpoints

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/users/create` | Create a user and hash its password |
| `POST` | `/users/login` | Verify credentials and issue an access cookie |
| `GET` | `/users/` | List users; requires a valid access cookie |

These are starter examples, not a complete production identity or authorization system. Review the [API reference](./docs/api-reference.md) and secure routes according to your application's requirements.

## Requirements and quick start

- Go `1.26.5` or newer.
- No external database server is required.

Copy `.env.example` to `.env`, then apply the migration and start the application:

```sh
cp .env.example .env
GOREST_DB_MIGRATION=true go run .
GOREST_DB_MIGRATION=false go run .
```

By default the database is stored at `app/database/gorest.db`. The application loads `.env` when present and also supports environment-only startup. See [Getting started](./docs/getting-started.md) and [SQLite configuration](./docs/database.md) for details.

## Tests and coverage

```sh
make unit-test                 # no database or external services required
make test-integration          # SQLite migration, seeder, library, and HTTP tests
make test-hardening-integration # local hardening endpoint test
make integration-reports       # tagged integration suite and coverage reports
make coverage                  # unit coverage reports
```

Each SQLite integration test creates its own database under a temporary directory and uses the repository's migration/seeder assets. It does not read or overwrite `app/database/gorest.db`. Optional Redis checks run only when Redis is explicitly enabled; they may connect to an external Redis instance. See the [testing guide](./docs/testing.md) before enabling those checks.

GitHub Actions runs unit and SQLite integration tests without database-server secrets. Integration coverage and test reports are uploaded as artifacts.

## Framework and security

The [GOREST documentation](https://github.com/elsyahtech/gorest/tree/main/docs) is the framework reference. Do not commit credentials, private keys, or certificates. Redis remains optional and disabled by the starter's `.env.example`; the current user endpoints do not depend on it.

For private vulnerability reporting and the project's security policy, see [SECURITY.md](./SECURITY.md).
