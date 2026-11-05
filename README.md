# Gorest API Starter

A ready-to-run application starter built with the [GOREST framework](https://github.com/elsyahtech/gorest). Use it to start an application with GOREST already wired, learn the framework through working examples, and adapt the structure and configuration to your own service. It includes layered CRUD routes, ORM and driver-specific native query examples, migrations, seeders, authentication examples, and tests.

This repository supports adoption of GOREST by providing a complete, modifiable application around the framework. The starter documentation explains this repository's choices; for framework APIs and guarantees, see the [GOREST documentation](https://github.com/elsyahtech/gorest/tree/main/docs). Adapt and secure the examples for your business requirements before exposing a service.

## Documentation

Start with the [documentation index](./docs/README.md). The guides are ordered as a learning path:

1. [Getting started](./docs/getting-started.md) — local setup and first requests.
2. [API reference](./docs/api-reference.md) — routes, request conventions, and examples.
3. [Architecture and lifecycle](./docs/architecture.md) — packages, end-to-end request flow, and naming conventions.
4. [Database guide](./docs/database.md) — drivers, migrations, seeders, and transaction limits.
5. [Configuration and deployment](./docs/configuration.md) — server, database, logging, and security considerations.
6. [Adding features](./docs/extending.md) — a step-by-step guide to extend the starter consistently.
7. [Testing and E2E coverage](./docs/testing.md) — run the suite, generate coverage, and verify the seven database drivers.

These guides show how this starter uses GOREST; they do not replace the framework reference. For framework APIs and behavior, see the [GOREST documentation](https://github.com/elsyahtech/gorest/tree/main/docs).

For vulnerability reporting and the project's security policy, see [SECURITY.md](./SECURITY.md).

## Important: examples are not production-ready

Most routes under `/examples` are unauthenticated. `POST /examples/login` issues an access cookie, and only `GET /examples/orm-queries/users/with-authentication` requires that token; this demonstration does not provide authorization for the other routes or production-ready access control. Do not expose the examples to untrusted clients. Before adapting this starter for a business application, secure every route and follow the [deployment checklist](./docs/configuration.md).

The checked-in configuration is for local examples only: it uses GOREST's SQLite default and enables both migrations and sample seeders (`Migration=true`, `Seeder=true`). Do not deploy these startup settings to production. Production application instances must not run the sample seeders or apply migrations automatically; use a controlled deployment migration step. See the [database guide](./docs/database.md) before selecting another driver or adapting the configuration for deployment.

## Requirements

- Go `1.26.5` or newer.
- SQLite for the default configuration; no external database service is required to start the examples.

Migration and seeder examples are provided for PostgreSQL, MySQL, Oracle, SQL Server, SQLite, MongoDB, and ScyllaDB. Their schema, query syntax, and transaction capabilities are not identical. See the [database guide](./docs/database.md) before switching drivers.

## Quick start

1. Review the local server and CORS settings in `app/config/server.go`.
2. From the repository root, run:

   ```sh
   go run .
   ```

   Or run `make run`. The server listens at `http://127.0.0.1:3000`. Database migrations and seeders run at startup; the SQLite file is created at `app/database/gorest.db`.
3. Try a read-only example:

   ```sh
   curl http://127.0.0.1:3000/examples/orm-queries/users
   ```

For complete setup steps, payload examples, and endpoint behavior, use [Getting started](./docs/getting-started.md) and the [API reference](./docs/api-reference.md).

## Development checks

```sh
make test      # run all Go tests with per-test progress
make check     # run go vet and all Go tests with progress and summaries
make tool-versions # verify pinned Go, golangci-lint, and Gitleaks versions
make lint      # pinned baseline golangci-lint checks (install separately)
make lint-full # all configured linters, including the quality backlog
make build     # build bin/gorest-api-starter
make test-e2e          # run the complete isolated E2E matrix for all seven drivers
make test-e2e-sqlite   # run the matrix for one driver (same targets exist for all drivers)
make coverage     # generate SQLite E2E coverage.out and coverage.html
make coverage-e2e # generate combined native-driver coverage (services required)
make push          # secret-scan, then interactively sanitize, stage, commit, and optionally push
```

`make test` runs the standard Go suite without starting external database services. `make test-e2e` is the isolated database workflow: it runs the verification scenarios against disposable containers for all seven drivers, including native SQL routes for the five SQL drivers, native MongoDB routes, and native CQL routes for ScyllaDB. It requires Docker and uses additional time, disk, and memory, especially for Oracle. Driver-specific targets (`make test-e2e-postgres`, `make test-e2e-mysql`, `make test-e2e-oracle`, `make test-e2e-sqlserver`, `make test-e2e-mongo`, `make test-e2e-scylla`, and `make test-e2e-sqlite`) run one part of the same matrix. See the [testing guide](./docs/testing.md) for prerequisites and scenario coverage.

The automated database testing supports macOS and Linux. Windows is not currently in the verified testing scope.

See the [testing and E2E guide](./docs/testing.md) for environment variables, coverage details, and test limitations. A passing test or build does not establish that every example works with every supported database driver.

`make push` runs Gitleaks against the current working tree and Git history before sanitizing or staging. It then asks for confirmation before sanitizing, staging, committing, and pushing. Before staging, it temporarily blanks non-empty `Password`, `Secret`, `Token`, `APIKey`, `ClientSecret`, `AccessKey`, and `PrivateKey` string fields in `app/config/*.go`; it restores the working files after staging, so the commit contains the sanitized values while the local configuration remains unchanged. It stages all repository changes, requires a clean Git index at startup, and never force-pushes. Review the staged summary and your branch/remote before confirming. The sanitizer is a narrow helper; Gitleaks is the repository and history scanner.
