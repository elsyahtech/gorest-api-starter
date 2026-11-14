# Extending the user example

The starter's active example is the SQLite-backed user API. Keep route definitions in `app/router`, transport handling in `app/handler`, business behavior in `app/library`, persistence in `app/model`, and public response mapping in `app/view`.

When adding behavior:

1. Define or update request/entity and public response types in `app/entity`.
2. Add business rules to the library layer, keeping persistence behind the model.
3. Use a narrow interface where a unit test needs a repository or use-case stub.
4. Register the route in `app/router` and apply authentication/authorization deliberately.
5. Add isolated tests in the package that owns the behavior.
6. Add or extend HTTP scenarios in `integration-test` when a complete request-to-SQLite flow needs verification.
7. Update the API and testing documentation and regenerate coverage/report outputs.

SQLite integration tests should create isolated database files in temporary directories and use the repository's migration/seeder assets. Do not rely on or overwrite the default application database in tests. Optional Redis integration tests should continue to use dedicated non-production test infrastructure.
