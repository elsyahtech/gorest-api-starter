# Example API reference

All endpoints on this page use the `/examples` prefix. They demonstrate application features and do not have authentication or authorization. Do not use them as a public business API without securing and adapting the implementation.

Default local base URL: `http://127.0.0.1:3000`.

## Request conventions

- Request bodies use JSON. Include the `Content-Type: application/json` header.
- The client supplies IDs; the examples do not generate them automatically.
- User create requests require `id`, `email`, `password`, and `departmentId`. Use a department ID that exists in the selected database so the example relationship is valid; database constraints differ by driver. The application processes the password before storing it.
- A department has `id`, `name`, `description`, and `isActive` fields.
- An order has `id`, `userId`, `orderNumber`, `status`, `totalAmount`, and `isActive` fields. Use a `userId` that exists in the selected database so the example relationship is valid; database constraints differ by driver.
- User responses use a DTO that omits the password. The JSON envelope and status code follow the GOREST response helper used by the application.
- `PATCH` accepts only the fields supported by that resource. An omitted field is not treated as its zero value.
- Send the complete resource body for standard `PUT` routes. The selected-columns user route is a documented exception that updates only the columns described below. Use `PATCH` for the explicitly supported partial updates listed below.
- `HEAD` is registered for the corresponding read routes. Explicit `OPTIONS` routes are registered for the standard user, department, and order resources; the query-example and native-driver routes do not all have dedicated `OPTIONS` handlers.

### Request body examples

Create a department:

```json
{
  "id": "dept-local-01",
  "name": "Engineering",
  "description": "Local development department",
  "isActive": true
}
```

Create a user:

```json
{
  "id": "user-local-01",
  "email": "user@example.com",
  "password": "local-test-password",
  "departmentId": "dept-01",
  "isActive": true
}
```

Create an order:

```json
{
  "id": "order-local-01",
  "userId": "user-local-01",
  "orderNumber": "ORD-LOCAL-01",
  "status": "PENDING",
  "totalAmount": 25.5,
  "isActive": true
}
```

Bulk user operations accept a JSON array of user objects using the same fields as the single-user examples:

```json
[
  {
    "id": "user-local-02",
    "email": "user2@example.com",
    "password": "local-test-password",
    "departmentId": "dept-01",
    "isActive": true
  }
]
```

## Departments

| Method | Path | Behavior |
|---|---|---|
| `GET` | `/examples/departments` | List departments |
| `GET` | `/examples/departments/:id` | Get one department |
| `HEAD` | `/examples/departments` | Return collection metadata without a body |
| `HEAD` | `/examples/departments/:id` | Return resource metadata without a body |
| `OPTIONS` | `/examples/departments` | List methods available on the collection |
| `OPTIONS` | `/examples/departments/:id` | List methods available on the resource |
| `POST` | `/examples/departments` | Create a department |
| `PUT` | `/examples/departments/:id` | Update a department |
| `PATCH` | `/examples/departments/:id` | Update `name`, `description`, and/or `isActive` |
| `DELETE` | `/examples/departments/:id` | Delete a department |

Example patch body:

```json
{"isActive": false}
```

## Orders

| Method | Path | Behavior |
|---|---|---|
| `GET` | `/examples/orders` | List orders |
| `GET` | `/examples/orders/:id` | Get one order |
| `GET` | `/examples/users/:userId/orders` | List a user's orders |
| `HEAD` | Each of the three `GET` paths above | Return metadata without a body |
| `OPTIONS` | `/examples/orders` | List methods available on the collection |
| `OPTIONS` | `/examples/orders/:id` | List methods available on the resource |
| `OPTIONS` | `/examples/users/:userId/orders` | List methods available on a user's order collection |
| `POST` | `/examples/orders` | Create an order |
| `PUT` | `/examples/orders/:id` | Update an order |
| `PATCH` | `/examples/orders/:id` | Update `orderNumber`, `status`, `totalAmount`, and/or `isActive` |
| `DELETE` | `/examples/orders/:id` | Delete an order |

## Users with the ORM

### Read and query examples

| Method | Path | Behavior |
|---|---|---|
| `GET` | `/examples/orm-queries/users` | List users |
| `GET` | `/examples/orm-queries/users/:id` | Get a user by ID |
| `HEAD` | `/examples/orm-queries/users` | Return collection metadata without a body |
| `HEAD` | `/examples/orm-queries/users/:id` | Return resource metadata without a body |
| `OPTIONS` | `/examples/orm-queries/users` | List methods available on the collection |
| `OPTIONS` | `/examples/orm-queries/users/:id` | List methods available on the resource |
| `GET` | `/examples/orm-queries/users/select-all` | Select all columns in the query |
| `GET` | `/examples/orm-queries/users/select-columns` | Demonstrate selecting specific columns |
| `GET` | `/examples/orm-queries/users/order-by` | Sort by ID in descending order |
| `GET` | `/examples/orm-queries/users/group-by` | Count users per department |
| `GET` | `/examples/orm-queries/users/distinct` | Demonstrate a distinct query on user department IDs |
| `GET` | `/examples/orm-queries/users/having` | Find departments with more than two users |
| `GET` | `/examples/orm-queries/users/join` | Join users and departments |
| `GET` | `/examples/orm-queries/users/preload` | Preload the department relation |
| `GET` | `/examples/orm-queries/users/limit?limit=5` | Limit the number of results |
| `GET` | `/examples/orm-queries/users/pagination?limit=5&offset=0` | Demonstrate limit and offset |

The `limit` and `offset` values must be valid integers. The `limit` route requires `limit`; the `pagination` route requires both values.

### Write examples

| Method | Path | Behavior |
|---|---|---|
| `POST` | `/examples/orm-queries/users` | Create a user |
| `POST` | `/examples/orm-queries/users/transaction` | Create a user using the transaction example |
| `POST` | `/examples/orm-queries/users/selected-columns` | Demonstrate creating selected columns |
| `POST` | `/examples/orm-queries/users/upsert` | Demonstrate an upsert operation |
| `POST` | `/examples/orm-queries/users/returning` | Demonstrate creating a user and returning the result |
| `POST` | `/examples/orm-queries/users/bulk` | Create users in bulk; body is a JSON array |
| `PUT` | `/examples/orm-queries/users/:id` | Update a user |
| `PUT` | `/examples/orm-queries/users/transaction/:id` | Update a user using the transaction example |
| `PUT` | `/examples/orm-queries/users/selected-columns/:id` | Update selected columns |
| `PUT` | `/examples/orm-queries/users/returning/:id` | Demonstrate updating a user and returning the result |
| `PUT` | `/examples/orm-queries/users/bulk` | Update users in bulk; body is a JSON array |
| `PATCH` | `/examples/orm-queries/users/:id` | Update only `isActive` |
| `DELETE` | `/examples/orm-queries/users/:id` | Delete a user |
| `DELETE` | `/examples/orm-queries/users/transaction/:id` | Delete a user using the transaction example |

Example user patch body:

```json
{"isActive": false}
```

User create operations and update operations that write a password hash the password before storing it. The selected-columns update is an exception: it updates only `id`, `email`, and `updated_at`, so it does not change the password. To change only the active status, use `PATCH`; you do not need to send a password.

The `upsert` example uses email as its conflict key. Create/returning and update/returning examples demonstrate returning selected result data; they are useful as ORM examples, not a promise that every database driver supports identical returning behavior.

The ORM delete routes call GOREST's `Delete` operation, while the native routes call the respective driver's delete operation. None of these handlers explicitly sets `deleted_at`; the presence of that nullable column in a schema does not itself implement soft deletion.

## Users with native database drivers

| Method | Path | Demonstrated driver/API |
|---|---|---|
| `GET`, `POST` | `/examples/native-queries/users/sql` | Native SQL |
| `PUT`, `DELETE` | `/examples/native-queries/users/sql/:id` | Native SQL |
| `GET`, `POST` | `/examples/native-queries/users/mongodb` | Native MongoDB |
| `PUT`, `DELETE` | `/examples/native-queries/users/mongodb/:id` | Native MongoDB |
| `GET`, `POST` | `/examples/native-queries/users/scylla` | Native ScyllaDB/CQL |
| `PUT`, `DELETE` | `/examples/native-queries/users/scylla/:id` | Native ScyllaDB/CQL |

Use the native routes that match your active database driver. The native SQL update changes email, password, department, and `updatedAt` only when the stored user is active; it does not change `isActive`. The native MongoDB update changes email, password, department, and `updatedAt`; it also leaves `isActive` unchanged. The native ScyllaDB update changes email, password, department, and `isActive`. The SQL examples use PostgreSQL syntax and `$1` placeholders, so they are not automatically portable to other SQL drivers. The MongoDB and ScyllaDB operations use those drivers' native APIs. GOREST transaction calls in the MongoDB and ScyllaDB ORM examples do not provide SQL-style transaction atomicity. See the [database guide](./database.md).

For the E2E coverage by route, local test commands, and opt-in database requirements, see [Testing and E2E coverage](./testing.md).
