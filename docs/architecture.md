# Architecture

This starter demonstrates a ScyllaDB-backed user API using small, focused application layers. The structure is a reference for clients adapting GOREST; it does not require every application to add layers it does not need.

## Request flow

```text
HTTP request
  -> app/router
  -> app/handler
  -> app/library
  -> app/model
  -> GOREST ScyllaDB database service
  -> ScyllaDB keyspace and users table
```

Responses return through the handler. `app/view` maps user entities to public response DTOs so password hashes are not included in user-list responses.

## Package responsibilities

| Package | Responsibility |
|---|---|
| `app/config` | Loads application, server, ScyllaDB, optional Redis, filter, logging, and monitoring settings. |
| `app/router` | Registers the active `/users/` routes and applies authentication to user listing. |
| `app/handler` | Parses HTTP requests, calls the application layer, and builds response envelopes. |
| `app/library` | Holds user application behavior, including password hashing and login verification. |
| `app/model` | Performs user data access through GOREST's ScyllaDB API. |
| `app/entity` | Defines the user entity and public response data structures. |
| `app/view` | Converts entities into public response DTOs. |
| `app/database/migrations/scylla` | Defines the users table and email secondary index. |
| `app/database/seeders/scylla` | Contains optional sample user rows. |
| `integration-test` | Verifies ScyllaDB connectivity, migrations, library behavior, HTTP behavior, and optional Redis behavior. |

There is no separate service package: `app/library` owns the current application logic. Redis is optional infrastructure wiring, disabled by default; the user endpoints do not depend on it.

## Active endpoints

- `POST /users/create` decodes and validates a user request, hashes the password, and persists the user.
- `POST /users/login` looks up the user, verifies the password, and issues an access cookie.
- `GET /users/` requires authentication, reads users from ScyllaDB, and returns public DTOs.

Creating users and logging in are public in this example. Authentication on the list route does not provide role-based authorization.

## ScyllaDB behavior

The users table uses `id` as its partition key and an email secondary index. Email lookup uses `ALLOW FILTERING` to include the active-user predicate; listing users scans the table. These patterns are suitable only for this small starter example. Design partition keys and query tables around production access patterns and data volume.

ScyllaDB does not provide a general transaction through GOREST's `BeginTx()` API. The user create path therefore performs a single conditional insert directly; multi-statement workflows need a Scylla-specific consistency design.

## Extending the example

Keep HTTP concerns in handlers, application behavior in the library, and persistence in the model. Add unit tests for logic that can be isolated without external services. Add ScyllaDB integration tests when behavior depends on real persistence, migrations, CQL, or TLS. Use a dedicated confirmed non-production keyspace, unique test data, and targeted cleanup.

See the [testing guide](./testing.md) for test commands, database safeguards, and coverage reports.
