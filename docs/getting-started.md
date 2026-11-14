# Getting started

This guide starts the SQL Server user API locally or against a SQL Server instance that you control.

## Requirements

- Go `1.26.5` or newer.
- A SQL Server instance reachable from the application host.
- An existing SQL Server database and login with the privileges needed to create and access the `users` table.

## Configure the application

Copy the example environment file and edit the SQL Server connection settings, or provide the settings directly in the process environment:

```sh
cp .env.example .env
```

Set `GOREST_DB_DRIVER=sqlserver`, `GOREST_DB_HOST`, `GOREST_DB_PORT` (default `1433`), `GOREST_DB_NAME`, `GOREST_DB_USER`, and `GOREST_DB_PASSWORD`. The example configures mandatory TLS and server-certificate validation; ensure the SQL Server certificate chains to a root trusted by the host and matches the configured hostname. Review [TLS and server certificate validation](./database.md#tls-and-server-certificate-validation). Do not put real credentials in source control.

The application loads `.env` when present; it can also start with environment variables alone. `main.go` applies the configured SQL Server migration when `GOREST_DB_MIGRATION=true`. Apply the migration once, then disable automatic migration for normal application runs:

```sh
GOREST_DB_MIGRATION=true go run .
GOREST_DB_MIGRATION=false go run .
```

The migration creates the SQL Server `dbo.users` table if it does not exist. The sample seeder is disabled by default; when enabled, it adds demonstration rows with placeholder password hashes, not accounts you can log in with. Never run the down migration against data you need: it drops the `users` table.

The server listens at the configured address, locally `http://127.0.0.1:3000`. To list users, first create a user and log in; see the [API reference](./api-reference.md).

## Redis

Redis is optional and is disabled by default in `.env.example`; the current user endpoints do not use it. If your application needs Redis, enable it and configure the connection and authentication settings. For TLS with a system-trusted server certificate, `GOREST_REDIS_TLS=true` is sufficient; a private CA or mTLS setup uses the certificate paths described in the [Redis configuration guide](./configuration.md#redis-optional). GOREST verifies connectivity with `PING` during startup.

## Run automated tests

```sh
make unit-test
make coverage
```

These commands do not connect to SQL Server. For the opt-in end-to-end suite, configure a dedicated non-production SQL Server database and exact-name confirmation as described in the [testing guide](./testing.md), then run `make test-integration`.
