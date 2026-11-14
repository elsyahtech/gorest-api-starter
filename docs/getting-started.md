# Getting started

This guide starts the ScyllaDB user API locally or against a ScyllaDB server that you control.

## Requirements

- Go `1.26.5` or newer.
- A ScyllaDB server reachable from the application host.
- A pre-created ScyllaDB keyspace and an account with the privileges needed by the application and migration.

## Configure the application

Copy the example environment file and edit the ScyllaDB connection settings, or provide the settings directly in the process environment:

```sh
cp .env.example .env
```

Set `GOREST_DB_HOST`, `GOREST_DB_PORT`, and `GOREST_DB_NAME` to the server and existing keyspace. `GOREST_DB_USER` and `GOREST_DB_PASSWORD` are optional but must be configured together. Review [TLS and client certificates](./database.md#tls-and-client-certificates) for server verification, private CA setup, and mTLS. Do not put real credentials or private keys in source control.

The application loads `.env` when present; it can also start with environment variables alone. `main.go` applies the configured ScyllaDB migration when `GOREST_DB_MIGRATION=true`. Run the migration after creating the keyspace, then disable automatic migration for normal application runs:

```sh
GOREST_DB_MIGRATION=true go run .
GOREST_DB_MIGRATION=false go run .
```

The migration creates the ScyllaDB `users` table and email index. The sample seeder is disabled by default. Never run the down migration against data you need: it drops the `users` table and index.

The server listens at the configured address, locally `http://127.0.0.1:3000`. To list users, first create a user and log in; see the [API reference](./api-reference.md).

## Redis

Redis is optional and is disabled by default in `.env.example`; the current user endpoints do not use it. If your application needs Redis, enable it and configure the connection and authentication settings. For TLS with a system-trusted server certificate, `GOREST_REDIS_TLS=true` is sufficient; a private CA or mTLS setup uses the certificate paths described in the [Redis configuration guide](./configuration.md#redis-optional). GOREST verifies connectivity with `PING` during startup.

## Run automated tests

```sh
make unit-test
make coverage
```

These commands do not connect to ScyllaDB. For the opt-in integration suite, configure a dedicated external test keyspace with exact-name confirmation as described in the [testing guide](./testing.md), then run `make test-integration`.
