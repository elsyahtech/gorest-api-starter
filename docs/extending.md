# Adding features

This guide explains the patterns used in the starter application. For route syntax, helpers, ORM usage, and framework APIs, see the [GOREST v1.0.2 documentation](https://github.com/elsyahtech/gorest/tree/v1.0.2/docs).

## Change flow

To add an operation to a resource:

1. **Entity** — Add or update the request or domain struct in `app/entity`. Set JSON and database mapping tags as needed. For `PATCH`, use pointers for optional fields so the application can recognize a supplied zero value.
2. **Model** — Add the query in `app/model`. Keep database access here, and use safe parameter binding for native queries.
3. **Library** — Add application rules and coordination in `app/library`, including input checks and any transformations needed before storing data.
4. **Handler** — In `app/handler`, parse parameters and request bodies, call the library, return errors using the existing GOREST response pattern, and log errors with the application logger.
5. **View** — If the response differs from the entity, add a DTO or mapping in `app/view`. Do not return sensitive fields such as passwords.
6. **Router** — Register the handler in the resource router under `app/router`. Make sure the HTTP method and path match the API.
7. **Database** — If the schema changes, update migrations and seeders for every driver the application supports.
8. **Documentation and checks** — Update the API reference and related guides. Run `make check` and add or update focused tests.

Follow the boundaries in [Architecture and application lifecycle](./architecture.md): handlers own HTTP input/output, libraries own application rules, and models own database operations. Keep the same operation name across layers where it describes the same behavior.

## Naming conventions

- Use lowercase filenames with underscores: `<resource>_<operation>_<layer>.go`, for example `department_read_handler.go`.
- Use the layer suffix (`_handler.go`, `_library.go`, or `_model.go`) so a file's responsibility is clear from its name.
- Name exported operations as **verb + resource + behavior**: `CreateUser`, `GetUserByID`, `UpdateUsersInBulk`, `PatchUserActiveStatus`.
- Prefer `Get` for reads and describe special behavior explicitly, such as `GetUsersWithNativeMongoQuery`.
- Keep names aligned across handler, library, model, and router. Avoid different terms for one operation unless the layers truly perform different behavior.
- Add a Go doc comment directly above every exported type, function, and method. Describe what it actually does and call out important limitations or filters.

The examples are grouped by resource, operation, and layer. Put package-level shared structs and constructors in `handler.go`, `library.go`, or `model.go`; reserve `interface.go` for actual Go interfaces.

## Registering routes

Each router currently creates a model, library, and handler for a resource, then registers routes under the `/examples` group. If you have both a static route and a parameter route, such as `/resource/search` and `/resource/:id`, register the static route first so it is not treated as an ID.

Add each new router to the `app.RegisterServices` list in `main.go`. If the endpoint supports different methods, parameters, or request bodies, also update its `HEAD` and `OPTIONS` routes where appropriate.

## Handle user updates carefully

The library calls the GOREST Argon helper before storing passwords in user create and update operations that write the password. The selected-columns update is an exception: it updates only `id`, `email`, and `updated_at`. Apply password hashing to any new operation that accepts and stores a password. Use response DTOs from `app/view` so passwords are not sent to clients.

`PATCH /examples/orm-queries/users/:id` currently supports only `isActive`. Department and order have different patch fields. If you add a patch field, update the patch request type, library checks, model update, API examples, and tests together.

## Transactions and multiple drivers

Follow the existing transaction pattern for operations that need transactions on databases that support them. The application includes migration and seeder files for several drivers, but SQL syntax, constraints, and driver capabilities can differ. Adapt native queries for each driver. Do not describe an operation as portable unless you have run and tested it with that driver.

## Done checklist

- The route is registered, and its method, path, request body, query parameters, response, and errors match the documentation.
- Input checks and error handling follow the repository patterns; failures are not ignored.
- Sensitive fields are not exposed, and native queries use suitable parameter binding.
- Migrations and sample data are consistent for the target driver.
- The method name and Go doc comment are consistent with the same operation in other layers.
- Relevant tests, `go vet`, and lint checks pass.
- The README and affected guides were reviewed for stale route or behavior descriptions.
