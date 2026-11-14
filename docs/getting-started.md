# Getting started

This guide starts the Oracle user API locally or against an Oracle instance that you control.

## Requirements

- Go `1.26.5` or newer.
- A compatible Oracle Instant Client.
- An existing Oracle service and application schema with privileges and tablespace quota to create the GOREST migration lock/history tables and the starter `users` table.

## Configure the application

Copy the example environment file and edit the Oracle connection settings:

```sh
cp .env.example .env
```

Set `GOREST_DB_DRIVER=oracle`, `GOREST_DB_HOST`, `GOREST_DB_PORT` (default `1521`), `GOREST_DB_NAME` (Oracle service name, not SID), `GOREST_DB_USER`, and `GOREST_DB_PASSWORD`. Set `GOREST_DB_LIB_DIR_PATH` if Oracle Instant Client is not discoverable through the system library path.

For Oracle TCPS, set `GOREST_DB_TLS=enable` and `GOREST_DB_CA_PATH` to the Oracle wallet directory used by `TNS_ADMIN`. The wallet and Oracle Net configuration must be compatible with the target listener and client. `GOREST_DB_TLS=disable` uses plain TCP; do not use it over untrusted networks. See [Oracle configuration](./database.md) for details.

The application loads `.env` when present; environment-only startup is also supported. Apply migrations on a controlled startup, then disable automatic migration for normal application runs:

```sh
GOREST_DB_MIGRATION=true go run .
GOREST_DB_MIGRATION=false go run .
```

The migration creates the `users` table and active-user index if needed. The sample seeder is disabled by default; when enabled, it adds demonstration rows with placeholder password hashes, not accounts you can log in with. Never run the down migrations against data you need.

The server listens at the configured address, locally `http://127.0.0.1:3000`. To list users, first create a user and log in; see the [API reference](./api-reference.md).

## Redis

Redis is optional and disabled by default in `.env.example`; the current user endpoints do not use it. If your application needs Redis, enable it and configure the connection and authentication settings. For TLS with a system-trusted server certificate, `GOREST_REDIS_TLS=true` is sufficient; a private CA or mTLS setup uses the certificate paths described in the [Redis configuration guide](./configuration.md#redis-optional). GOREST verifies connectivity with `PING` during startup.

## Run automated tests

```sh
make unit-test
make test-integration
make coverage
```

Oracle integration tests require the Instant Client, dedicated test schema, and exact schema confirmation described in the [testing guide](./testing.md).
