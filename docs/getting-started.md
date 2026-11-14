# Getting started

This guide starts the SQLite user API locally.

## Requirements

- Go `1.26.5` or newer.
- No external database server is required.

## Configure the application

Copy the example environment file:

```sh
cp .env.example .env
```

The example selects SQLite and stores the database at `app/database/gorest.db` by default. Change `GOREST_DB_NAME` or `GOREST_DB_DIRECTORY` if you want to use another local file location. The application loads `.env` when present and can also start with environment variables alone.

Apply the migration on a controlled startup, then disable automatic migration for normal runs:

```sh
GOREST_DB_MIGRATION=true go run .
GOREST_DB_MIGRATION=false go run .
```

The migration creates the `users` table if it does not exist. The sample seeder is disabled by default; when enabled, it adds demonstration rows with placeholder password hashes, not accounts you can log in with. Never run the down migration against data you need: it drops the `users` table.

The server listens at the configured address, locally `http://127.0.0.1:3000`. To list users, first create a user and log in; see the [API reference](./api-reference.md).

## Redis

Redis is optional and disabled by default in `.env.example`; the current user endpoints do not use it. If your application needs Redis, enable it and configure the connection and authentication settings. For TLS with a system-trusted server certificate, `GOREST_REDIS_TLS=true` is sufficient; a private CA or mTLS setup uses the certificate paths described in the [Redis configuration guide](./configuration.md#redis-optional). GOREST verifies connectivity with `PING` during startup.

## Run automated tests

```sh
make unit-test
make test-integration
make coverage
```

The integration suite uses isolated temporary SQLite files and does not modify the default application database. See the [testing guide](./testing.md) for the full suite and optional Redis checks.
