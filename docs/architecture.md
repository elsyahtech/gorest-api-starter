# Architecture

This starter demonstrates a MongoDB-backed user API using small, focused application layers. The structure is a reference for clients adapting GOREST; it does not require every application to add layers it does not need.

## Request flow

```text
HTTP request
  -> app/router
  -> app/handler
  -> app/library
  -> app/model
  -> GOREST MongoDB database service
  -> MongoDB database and users collection
```

Responses return through the handler. `app/view` maps user entities to public response DTOs so password hashes are not included in user-list responses.

## Package responsibilities

| Package | Responsibility |
|---|---|
| `app/config` | Loads application, server, MongoDB, optional Redis, filter, logging, and monitoring settings. |
| `app/router` | Registers the active `/users/` routes and applies authentication to user listing. |
| `app/handler` | Parses HTTP requests, calls the application layer, and builds response envelopes. |
| `app/library` | Holds user application behavior, including password hashing and login verification. |
| `app/model` | Performs user data access through GOREST's MongoDB API. |
| `app/entity` | Defines the user entity and public response data structures. |
| `app/view` | Converts entities into public response DTOs. |
| `app/database/migrations/mongo` | Defines the users collection validator and indexes. |
| `app/database/seeders/mongo` | Contains optional sample user documents. |
| `integration-test` | Verifies MongoDB connectivity, migrations, library behavior, HTTP behavior, and optional Redis behavior. |

There is no separate service package: `app/library` owns the current application logic. Redis is optional infrastructure wiring, disabled by default; the user endpoints do not depend on it.

## Active endpoints

- `POST /users/create` decodes and validates a user request, hashes the password, and persists the document.
- `POST /users/login` looks up the user, verifies the password, and issues an access cookie.
- `GET /users/` requires authentication, reads users from MongoDB, and returns public DTOs.

Creating users and logging in are public in this example. Authentication on the list route does not provide role-based authorization.

## MongoDB behavior

The user create path uses the ORM directly rather than opening a MongoDB transaction. A single-document insert is atomic and works with standalone servers; transactions across multiple documents require a replica set or sharded cluster. The migration creates unique indexes for the application `id` and `email` fields.

## Extending the example

Keep HTTP concerns in handlers, application behavior in the library, and persistence in the model. Add unit tests for logic that can be isolated without external services. Add MongoDB integration tests when behavior depends on real persistence, migrations, indexes, or driver configuration. Use a dedicated confirmed non-production database, unique test data, and targeted cleanup.

See the [testing guide](./testing.md) for connection prerequisites, database safeguards, test commands, and coverage reports.
