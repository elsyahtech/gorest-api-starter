# Gorest API Starter Documentation

This documentation is the application-specific guide for developers adopting or adapting the starter as their application foundation. It explains how this repository wires GOREST, and what the included examples do. It does not replace the framework reference: use the [GOREST v1.1.15 documentation](https://github.com/elsyahtech/gorest/tree/v1.1.15/docs) for framework APIs and guarantees.

## Learning path

1. [Getting started](./getting-started.md) — prerequisites, local configuration, startup, and first requests.
2. [Example API reference](./api-reference.md) — routes, payloads, query parameters, and behavior.
3. [Architecture and application lifecycle](./architecture.md) — end-to-end flow, package responsibilities, and code conventions.
4. [Database, migrations, and seeders](./database.md) — configured drivers, data files, startup behavior, and limitations.
5. [Configuration, security, and deployment](./configuration.md) — server, database, logging, Redis, and deployment checklist.
6. [Adding features](./extending.md) — a consistent workflow for extending the application.
7. [Testing and E2E coverage](./testing.md) — local test commands, coverage, optional database services, and suite limitations.
8. [Release checks](./release-checks.md) — pinned tool versions, CI release gate, and evidence retention.

## Scope

- These guides describe the configuration and examples in this repository. For framework APIs, ORM behavior, middleware, and helpers, use the GOREST documentation.
- Routes under `/examples` are demonstrations, not a production API contract. Most are unauthenticated; `POST /examples/login` issues an access cookie and only `GET /examples/orm-queries/users/with-authentication` requires it. This is not production-ready authorization or access control.
- The checked-in database configuration uses SQLite and enables this starter's migrations and seeders, so the examples can run locally without a database server or credentials. See the [database guide](./database.md) before switching drivers; review startup and secret handling before deployment.
- Database examples are driver-specific. A migration or query example for one database is not a guarantee of equivalent behavior on another driver.
- Keep this documentation in sync when changing routes, request/response DTOs, configuration, migrations, seeders, or application behavior.
- For private vulnerability reporting and the project's security policy, see [`SECURITY.md`](../SECURITY.md).
