# Architecture

This starter demonstrates a PostgreSQL-backed user API using small, focused application layers. The structure is a reference for clients adapting GOREST; it does not require every application to add layers it does not need.

## Request flow

```text
HTTP request
  -> app/router
  -> app/handler
  -> app/library
  -> app/model
  -> GOREST PostgreSQL database service
  -> PostgreSQL
```

Responses return through the handler. `app/view` maps user entities to public response DTOs so password hashes are not included in user-list responses.

## Package responsibilities

| Package | Responsibility |
|---|---|
| `app/config` | Loads application, server, PostgreSQL, optional Redis, filter, logging, and monitoring settings. |
| `app/router` | Registers the active `/users/` routes and applies authentication to user listing. |
| `app/handler` | Parses HTTP requests, calls the application layer, and builds response envelopes. |
| `app/library` | Holds user application behavior, including password hashing and login verification. |
| `app/model` | Performs user data access through GOREST's PostgreSQL API. |
| `app/entity` | Defines the user entity and public response data structures. |
| `app/view` | Converts entities into public response DTOs. |
| `app/database/migrations/postgres` | Contains the PostgreSQL schema migration. |
| `app/database/seeders/postgres` | Contains sample PostgreSQL user data. |
| `integration-test` | Verifies PostgreSQL mTLS, migrations, seeders, HTTP behavior, and optional Redis service behavior. |

There is no separate service package: `app/library` owns the current application logic. Redis is optional infrastructure wiring, disabled by default; the user endpoints do not depend on it.

## Active endpoints

- `POST /users/create` decodes and validates a user request, hashes the password, and persists the user.
- `POST /users/login` looks up the user, verifies the password, and issues an access cookie.
- `GET /users/` requires authentication, reads users from PostgreSQL, and returns public DTOs.

Creating users and logging in are public in this example. Authentication on the list route does not provide role-based authorization.

## Extending the example

Keep HTTP concerns in handlers, application behavior in the library, and persistence in the model. Add unit tests for logic that can be isolated without external services. Add PostgreSQL integration tests when behavior depends on real persistence, migrations, or database configuration; use a dedicated non-production database, unique test data, and targeted cleanup.

See the [testing guide](./testing.md) for test commands, database safeguards, and coverage reports.
