# GOREST Oracle API Starter

An Oracle-focused application starter built with [GOREST](https://github.com/elsyahtech/gorest/tree/main). It demonstrates a small user API with authentication, password hashing, Oracle persistence, automated unit tests, and opt-in integration tests against a dedicated Oracle schema.

## Documentation

Start with the [documentation index](./docs/README.md). The guides cover setup, the active API, architecture, Oracle configuration, and testing.

## Included endpoints

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/users/create` | Create a user and hash its password |
| `POST` | `/users/login` | Verify credentials and issue an access cookie |
| `GET` | `/users/` | List users; requires a valid access cookie |

These are starter examples, not a complete production identity or authorization system. Review the [API reference](./docs/api-reference.md) and secure routes according to your application's requirements.

## Requirements and quick start

- Go `1.26.5` or newer.
- Oracle Instant Client, with `GOREST_DB_LIB_DIR_PATH` set if it is not in the system library path.
- An Oracle service name and application schema.

Copy `.env.example` to `.env`, configure the Oracle connection, then apply the migration and start the application:

```sh
cp .env.example .env
# Configure the Oracle service, schema, credentials, and client library.
GOREST_DB_MIGRATION=true go run .
GOREST_DB_MIGRATION=false go run .
```

`GOREST_DB_NAME` is the Oracle service name, not a SID. For TCPS, set `GOREST_DB_TLS=enable` and configure `GOREST_DB_CA_PATH` to an Oracle wallet directory. See [Getting started](./docs/getting-started.md) and [Oracle configuration](./docs/database.md).

## Tests and coverage

```sh
make unit-test                  # no Oracle server required
make test-integration           # Oracle connection, migration, library, and HTTP tests
make test-oracle-connection     # connection-only Oracle integration test
make test-hardening-integration # local hardening endpoint test
make integration-reports        # tagged integration suite and coverage reports
make coverage                   # unit coverage reports
```

Oracle integration tests perform schema changes and require a dedicated non-production schema. Before running them, set `GOREST_TEST_ORACLE_CONFIRM_USER` to exactly match `GOREST_DB_USER`. See the [testing guide](./docs/testing.md). Optional Redis checks run only when Redis is enabled and may connect to an external Redis instance.

GitHub Actions runs unit tests and isolated Oracle integration-configuration checks. It does not connect to an Oracle server; the full Oracle tests need an Oracle client, wallet when using TCPS, and a dedicated test schema.

## Framework and security

The [GOREST documentation](https://github.com/elsyahtech/gorest/tree/main/docs) is the framework reference. Do not commit credentials, private keys, wallets, or certificates. Redis remains optional and disabled by the starter's `.env.example`; the current user endpoints do not depend on it.

For private vulnerability reporting and the project's security policy, see [SECURITY.md](./SECURITY.md).
