# Getting started

This guide starts the MongoDB user API locally or against a MongoDB instance that you control.

## Requirements

- Go `1.26.5` or newer.
- A MongoDB server reachable from the application host.
- A MongoDB database and an account with permissions to create/access the `users` collection and its indexes, if migrations are enabled.

## Configure the application

Copy the example environment file and edit the MongoDB connection settings:

```sh
cp .env.example .env
```

Set `GOREST_DB_DRIVER=mongo`, `GOREST_DB_HOST`, `GOREST_DB_PORT` (default `27017`), and `GOREST_DB_NAME`. If authentication is enabled, set both `GOREST_DB_USER` and `GOREST_DB_PASSWORD`, and set `GOREST_DB_AUTH_SOURCE` if the account is stored outside `admin`. This starter constructs a standard `mongodb://host:port/database` URI; MongoDB SRV connection strings are not supported by these fields.

For a remote deployment, set `GOREST_DB_TLS=true` and configure the server certificate trust. Use `GOREST_DB_CA_PATH` for a private CA. For mTLS, configure both `GOREST_DB_CERT_PATH` and `GOREST_DB_KEY_PATH`; the MongoDB server must require and trust the client certificate. `GOREST_DB_TLS=false` is for trusted local or isolated networks only. See [MongoDB configuration](./database.md).

Apply the collection/index migration on a controlled startup, then disable automatic migration for normal application runs:

```sh
GOREST_DB_MIGRATION=true go run .
GOREST_DB_MIGRATION=false go run .
```

The migration creates the `users` collection, document validator, and indexes if needed. The sample seeder is disabled by default; enabling it inserts three example documents with placeholder hashes, not accounts you can log in with. The down migration drops the entire `users` collection, so never use it where its data must be retained.

The server listens locally at `http://127.0.0.1:3000`. To list users, first create a user and log in; see the [API reference](./api-reference.md).

## Redis

Redis is optional and disabled by default in `.env.example`; the current user endpoints do not use it. If your application needs Redis, enable it and configure the connection and authentication settings. For TLS with a system-trusted server certificate, `GOREST_REDIS_TLS=true` is sufficient; a private CA or mTLS setup uses the certificate paths described in the [Redis configuration guide](./configuration.md#redis-optional). GOREST verifies connectivity with `PING` during startup.

## Run automated tests

```sh
make unit-test
make coverage
```

These commands do not connect to MongoDB. For the opt-in external integration suite, configure a dedicated test database and its exact confirmation as described in the [testing guide](./testing.md), then run `make test-integration`.
