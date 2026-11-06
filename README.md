# GOREST PostgreSQL API Starter

A PostgreSQL-focused application starter built with [GOREST](https://github.com/elsyahtech/gorest/tree/main). It demonstrates a small user API with authentication, password hashing, PostgreSQL persistence, automated unit tests, and opt-in integration tests against a separately confirmed PostgreSQL test database.

## Documentation

Start with the [documentation index](./docs/README.md). The guides cover setup, the active API, architecture, PostgreSQL configuration, and testing evidence.

## Included endpoints

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/users/create` | Create a user and hash its password |
| `POST` | `/users/login` | Verify credentials and issue an access cookie |
| `GET` | `/users/` | List users; requires a valid access cookie |

These are starter examples, not a complete production identity or authorization system. Review the [API reference](./docs/api-reference.md) and secure routes according to your application's requirements.

## Requirements and quick start

- Go `1.26.5` or newer.
- A reachable PostgreSQL server and a database for the application.

Copy `.env.example` to `.env`, provide the PostgreSQL database name, username, password, and appropriate TLS settings, then apply the PostgreSQL migration and start the application:

```sh
cp .env.example .env
# Edit .env with your PostgreSQL settings.
GOREST_DB_MIGRATION=true go run .
GOREST_DB_MIGRATION=false go run .
```

The application loads `.env` when present and also supports environment-only startup. See [Getting started](./docs/getting-started.md) and [PostgreSQL configuration](./docs/database.md) for the full setup.

## Tests and coverage

```sh
make unit-test       # isolated unit tests; no database required
make coverage        # unit coverage profile, text summary, and HTML report
make test-integration # PostgreSQL TLS/mTLS and HTTP integration tests
make integration-reports # one integration run: readable and JSON reports plus app coverage
```

Integration tests are opt-in and use the PostgreSQL `GOREST_DB_*` settings in `.env` (or environment variables), plus an explicit database-name confirmation. They require a dedicated non-production test database; the HTTP test runs its PostgreSQL migration and seeder and writes a uniquely identified user. The optional Redis integration test uses `GOREST_REDIS_*` only when Redis is explicitly enabled; point those settings at a non-production test instance. Tests do not use Docker or a disposable local database. Read the [testing guide](./docs/testing.md) before running them.

Coverage and integration outputs are generated under `coverage/`. `make integration-reports` runs the full tagged integration suite once and produces test output and application coverage; if Redis is enabled in the environment, its optional checks may also connect to the configured Redis instance. See the [testing guide](./docs/testing.md) before running integration commands.

GitHub Actions runs unit tests and coverage on pushes and pull requests. Its PostgreSQL mTLS integration job runs when the required external test database and certificate secrets are configured; see the [testing guide](./docs/testing.md#github-actions-ci).

## Framework and security

The [GOREST documentation](https://github.com/elsyahtech/gorest/tree/main/docs) is the framework reference. For PostgreSQL TLS, private-CA, and mTLS setup, follow the [PostgreSQL configuration guide](./docs/database.md#configure-tls-private-cas-and-mutual-tls) and ensure the server's client-certificate policy matches your configuration. Do not commit credentials, private keys, or certificates. Redis remains optional and disabled by the starter's `.env.example`; the current user endpoints do not depend on it.

For private vulnerability reporting and the project's security policy, see [SECURITY.md](./SECURITY.md).
