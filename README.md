# Gorest API Starter

A Go REST API reference application built with [GOREST](https://github.com/elsyahtech/gorest). It demonstrates a layered application structure, CRUD endpoints, ORM and native database queries, transactions, migrations, seeders, and request/response DTOs.

This repository is intended as a **learning and implementation reference**. Adapt its examples to your business requirements before using them in an internal or public production service.

## Documentation

Start with the [documentation index](./docs/README.md). The guides are ordered as a learning path:

1. [Getting started](./docs/getting-started.md) — local setup and first requests.
2. [API reference](./docs/api-reference.md) — routes, request conventions, and examples.
3. [Architecture and lifecycle](./docs/architecture.md) — packages, end-to-end request flow, and naming conventions.
4. [Database guide](./docs/database.md) — drivers, migrations, seeders, and transaction limits.
5. [Configuration and deployment](./docs/configuration.md) — server, database, logging, and security considerations.
6. [Adding features](./docs/extending.md) — a step-by-step guide to extend the starter consistently.
7. [Testing and E2E coverage](./docs/testing.md) — run the suite, generate coverage, and enable optional database tests.

These guides describe this application, not every framework feature. For GOREST framework APIs and behavior, see the [GOREST v1.0.2 documentation](https://github.com/elsyahtech/gorest/tree/v1.0.2/docs).

For vulnerability reporting and the project's security policy, see [SECURITY.md](./SECURITY.md).

## Important: examples are not production-ready

All routes are registered under `/examples` and have no authentication or authorization. Do not expose them to untrusted clients. Before adapting this starter for a business application, add access controls, review every example operation, and follow the [deployment checklist](./docs/configuration.md).

The checked-in configuration is for local examples: PostgreSQL is enabled, a sample database password is present in source, and migrations and seeders run at startup. Do not use these values or startup settings in production.

## Requirements

- Go `1.26.5` or newer.
- PostgreSQL for the default configuration.

Migration and seeder examples are provided for PostgreSQL, MySQL, Oracle, SQL Server, SQLite, MongoDB, and ScyllaDB. Their schema, query syntax, and transaction capabilities are not identical. See the [database guide](./docs/database.md) before switching drivers.

## Quick start

1. Create a local PostgreSQL database and user.
2. Update the local connection settings in `app/config/database.go`. The checked-in credentials are examples.
3. Review the local server and CORS settings in `app/config/server.go`.
4. From the repository root, run:

   ```sh
   go run .
   ```

   Or run `make run`. The server listens at `http://127.0.0.1:3000`. Database migrations and seeders are enabled by default and run during application startup.
5. Try a read-only example:

   ```sh
   curl http://127.0.0.1:3000/examples/orm-queries/users
   ```

For complete setup steps, payload examples, and endpoint behavior, use [Getting started](./docs/getting-started.md) and the [API reference](./docs/api-reference.md).

## Development checks

```sh
make test      # go test ./...
make check     # go vet ./... and go test ./...
make lint      # golangci-lint run ./... (install separately)
make build     # build bin/gorest-api-starter
make test-e2e  # run E2E tests; external database tests are skipped unless enabled
make test-mongo   # run opt-in MongoDB tests (requires a local service)
make test-scylla  # run opt-in ScyllaDB tests (requires a local service)
make coverage     # generate SQLite E2E coverage.out and coverage.html
make coverage-e2e # include opt-in MongoDB and ScyllaDB tests
make push          # interactively sanitize, stage, commit, and optionally push
```

See the [testing and E2E guide](./docs/testing.md) for service setup, environment variables, coverage details, and test limitations. The default suite includes SQLite-backed HTTP E2E tests; MongoDB and ScyllaDB tests are opt-in and require separately provisioned services. A passing test or build does not establish that every example works with every supported database driver.

`make push` asks for confirmation before sanitizing, staging, committing, and pushing. Before staging, it temporarily blanks non-empty `Password`, `Secret`, `Token`, `APIKey`, `ClientSecret`, `AccessKey`, and `PrivateKey` string fields in `app/config/*.go`; it restores the working files after staging, so the commit contains the sanitized values while the local configuration remains unchanged. It stages all repository changes, requires a clean Git index at startup, and never force-pushes. Review the staged summary and your branch/remote before confirming; this helper only sanitizes the listed field names in `app/config/*.go` and is not a general secret scanner.
