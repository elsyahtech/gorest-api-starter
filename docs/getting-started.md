# Getting started

This guide starts the MySQL user API locally or against a MySQL server that you control.

## Requirements

- Go `1.26.5` or newer.
- A MySQL server reachable from the application host.
- A MySQL database, user, and password with the privileges needed to create and access the `users` table.

## Configure the application

Copy the example environment file and edit the MySQL connection settings, or provide the settings directly in the process environment:

```sh
cp .env.example .env
```

Set `GOREST_DB_HOST`, `GOREST_DB_PORT`, `GOREST_DB_NAME`, `GOREST_DB_USER`, and `GOREST_DB_PASSWORD`. Before configuring TLS, review [TLS, private CAs, and client certificates](./database.md#configure-tls-private-cas-and-client-certificates) to learn how the connector uses system trust roots and how to add a private CA to the application process trust bundle. The current connector does not provide mTLS through the certificate/key path settings. Do not put real credentials in source control.

The application loads `.env` when present; it can also start with environment variables alone. `main.go` applies the configured MySQL migration when `GOREST_DB_MIGRATION=true`. Apply the migration once, then disable automatic migration for normal application runs:

```sh
GOREST_DB_MIGRATION=true go run .
GOREST_DB_MIGRATION=false go run .
```

The migration creates the MySQL `users` table if it does not exist. The sample seeder is disabled by default. Never run the down migration against data you need: it drops the `users` table.

The server listens at the configured address, locally `http://127.0.0.1:3000`. To list users, first create a user and log in; see the [API reference](./api-reference.md).

## Redis

Redis is optional and is disabled by default in `.env.example`. The current user endpoints do not use Redis. Set `GOREST_REDIS_DISABLE=false` only when your application has a Redis instance configured and needs the service.

## Run automated tests

```sh
make unit-test
make coverage
```

These commands do not connect to MySQL. For the opt-in end-to-end suite, configure a dedicated external MySQL test database as described in the [testing guide](./testing.md), then run `make test-integration`.
