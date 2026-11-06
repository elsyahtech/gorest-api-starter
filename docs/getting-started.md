# Getting started

This guide starts the PostgreSQL user API locally or against a PostgreSQL server that you control.

## Requirements

- Go `1.26.5` or newer.
- A PostgreSQL server reachable from the application host.
- A PostgreSQL database, user, and password with the privileges needed to create and access the `users` table.

## Configure the application

Copy the example environment file and edit the PostgreSQL connection settings, or provide the settings directly in the process environment:

```sh
cp .env.example .env
```

Set `GOREST_DB_DRIVER=postgres`, `GOREST_DB_HOST`, `GOREST_DB_PORT`, `GOREST_DB_NAME`, `GOREST_DB_USER`, and `GOREST_DB_PASSWORD`. Configure `GOREST_DB_TLS=verify-full` for server certificate and hostname verification. Review [TLS, private CAs, and mTLS](./database.md#configure-tls-private-cas-and-mutual-tls) for CA setup and optional client-certificate authentication. Do not put real credentials or private keys in source control.

The application loads `.env` when present; it can also start with environment variables alone. `main.go` applies the configured PostgreSQL migration when `GOREST_DB_MIGRATION=true`. Apply the migration once, then disable automatic migration for normal application runs:

```sh
GOREST_DB_MIGRATION=true go run .
GOREST_DB_MIGRATION=false go run .
```

The migration creates the PostgreSQL `users` table if it does not exist. The sample seeder is disabled by default; when enabled, it adds demonstration rows with placeholder password hashes, not accounts you can log in with. Never run the down migration against data you need: it drops the `users` table.

The server listens at the configured address, locally `http://127.0.0.1:3000`. To list users, first create a user and log in; see the [API reference](./api-reference.md).

## Redis

Redis is optional and is disabled by default in `.env.example`; the current user endpoints do not use it. If your application needs Redis, enable it and configure the connection and authentication settings. For TLS with a system-trusted server certificate, `GOREST_REDIS_TLS=true` is sufficient; a private CA or mTLS setup uses the certificate paths described in the [Redis configuration guide](./configuration.md#redis-optional). GOREST verifies connectivity with `PING` during startup.

## Run automated tests

```sh
make unit-test
make coverage
```

These commands do not connect to PostgreSQL. For the opt-in end-to-end suite, configure a dedicated non-production PostgreSQL database and mTLS credentials as described in the [testing guide](./testing.md), then run `make test-integration`.
