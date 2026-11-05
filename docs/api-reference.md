# User API reference

Base URL: `http://127.0.0.1:3000` unless changed in the server configuration. The application currently registers only the user routes below.

## Response envelope

Responses are JSON objects with `success`, `httpCode`, and either `message` or `error`. User-list data is transformed to a public response shape that omits password hashes.

## Create a user

`POST /users/create`

```json
{
  "id": "client-generated-user-id",
  "email": "user@example.test",
  "password": "a-test-password",
  "isActive": true
}
```

The client supplies the ID. The application hashes the password before inserting the user. Successful creation returns HTTP `200`; malformed or invalid request bodies return an error response.

## Log in

`POST /users/login`

```json
{
  "email": "user@example.test",
  "password": "a-test-password"
}
```

Valid credentials return HTTP `200` and set an `access` cookie. Invalid credentials return HTTP `401`. Configure a non-placeholder, high-entropy `GOREST_FILTER_SIGNATURE_KEY` before exposing the service.

## List users

`GET /users/`

This route requires a valid `access` cookie from the login endpoint. A request without valid authentication returns HTTP `401`. An authenticated request returns an array of public user DTOs containing ID, email, active state, and timestamps; it does not return the password hash.

## Authentication and scope

Only the list route uses the GOREST authentication filter. Creating users and logging in are public by design in this example. The filter demonstrates authentication, not role-based or resource-level authorization. Add the controls required by your application before exposing business data.

For automated success and error-path checks, see [Testing and coverage](./testing.md).
