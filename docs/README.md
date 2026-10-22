# Gorest API Starter Documentation

This documentation is the application-specific guide for developers evaluating or adapting the starter. Follow it in order for a guided introduction, or use the individual references as needed. It describes the code and examples in this repository; it does not replace the [GOREST v1.0.2 framework documentation](https://github.com/elsyahtech/gorest/tree/v1.0.2/docs).

## Learning path

1. [Getting started](./getting-started.md) — prerequisites, local configuration, startup, and first requests.
2. [Example API reference](./api-reference.md) — routes, payloads, query parameters, and behavior.
3. [Architecture and application lifecycle](./architecture.md) — end-to-end flow, package responsibilities, and code conventions.
4. [Database, migrations, and seeders](./database.md) — configured drivers, data files, startup behavior, and limitations.
5. [Configuration, security, and deployment](./configuration.md) — server, database, logging, Redis, and deployment checklist.
6. [Adding features](./extending.md) — a consistent workflow for extending the application.
7. [Testing and E2E coverage](./testing.md) — local test commands, coverage, optional database services, and suite limitations.

## Scope

- These guides describe the configuration and examples in this repository. For framework APIs, ORM behavior, middleware, and helpers, use the GOREST documentation.
- Routes under `/examples` are demonstrations, not a production API contract. They have no authentication or authorization.
- The checked-in configuration is for local development and includes a sample database password, enabled migrations, and enabled seeders. Replace these settings before deployment.
- Database examples are driver-specific. A migration or query example for one database is not a guarantee of equivalent behavior on another driver.
- Keep this documentation in sync when changing routes, request/response DTOs, configuration, migrations, seeders, or application behavior.
- For private vulnerability reporting and the project's security policy, see [`SECURITY.md`](../SECURITY.md).
