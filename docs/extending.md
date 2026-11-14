# Extending the user example

The starter's active example is the Oracle-backed user API. Keep route definitions in `app/router`, transport handling in `app/handler`, business behavior in `app/library`, persistence in `app/model`, and public response mapping in `app/view`.

When adding behavior:

1. Define or update request/entity and public response types in `app/entity`.
2. Add business rules to the library layer, keeping persistence behind the model.
3. Use a narrow interface where a unit test needs a repository or use-case stub.
4. Register the route in `app/router` and apply authentication/authorization deliberately.
5. Add isolated tests in the package that owns the behavior.
6. Add or extend HTTP scenarios in `integration-test` when a complete request-to-Oracle flow needs verification.
7. Update the API and testing documentation and regenerate coverage/report outputs.

Oracle integration tests require a dedicated non-production schema, unique test data, exact schema-user confirmation, an Oracle Instant Client, and targeted cleanup. Use a wallet when testing TCPS. Do not copy credentials into test code or use production data.
