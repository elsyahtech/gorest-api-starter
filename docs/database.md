# MongoDB configuration

This starter uses MongoDB through GOREST's database service. Runtime settings are read by `app/config/database.go`.

## Application connection

Set these variables in `.env`:

| Variable | Purpose |
|---|---|
| `GOREST_DB_DRIVER` | Must be `mongo` for this starter. |
| `GOREST_DB_HOST`, `GOREST_DB_PORT` | MongoDB server host and port (default `27017`). This starter builds a `mongodb://` host-and-port URI; it does not accept `mongodb+srv://` connection strings. |
| `GOREST_DB_NAME` | MongoDB database name; required. |
| `GOREST_DB_USER`, `GOREST_DB_PASSWORD` | Optional authentication credentials; configure both or neither. |
| `GOREST_DB_AUTH_SOURCE` | Database that stores the user's credentials; defaults to `admin`. |
| `GOREST_DB_TLS` | `true` enables TLS; `false` uses plaintext TCP. |
| `GOREST_DB_CA_PATH` | Optional CA certificate file for server verification. |
| `GOREST_DB_CERT_PATH`, `GOREST_DB_KEY_PATH` | Optional client certificate and key for mTLS; configure both together and enable TLS. |
| `GOREST_DB_TIMEOUT` | MongoDB server-selection and connection timeout; defaults to GOREST's configured timeout. |
| `GOREST_DB_MIGRATION` | Enables MongoDB collection and index migrations during application startup. |
| `GOREST_DB_SEEDER` | Enables sample document insertion; disabled by default. |

For remote deployments, set `GOREST_DB_TLS=true`. Use a trusted certificate and hostname; configure `GOREST_DB_CA_PATH` when the MongoDB server certificate is signed by a private CA. For mTLS, provide both the client certificate and matching private key, and ensure the server requires and trusts that identity. Keep certificate verification enabled and protect private keys.

The up migration in `app/database/migrations/mongo/` creates the `users` collection, validates its document fields, and adds unique indexes for `id` and `email` plus an index for `is_active`. The down migration drops the collection and its data; never run it against data that must be retained. The sample seeder is disabled by default and inserts three rows with placeholder password hashes, not usable login credentials. Its compensation file deletes only those fixed sample IDs.

The user model uses MongoDB document operations without wrapping each create in a transaction. This supports standalone MongoDB deployments; MongoDB multi-document transactions require a replica set or sharded cluster. See the [GOREST MongoDB guide](https://github.com/elsyahtech/gorest/blob/v1.1.21/docs/testing-result/testing-result-mongo.md) for framework-specific behavior and deployment considerations.

## Integration-test database

Integration tests accept explicit `GOREST_TEST_MONGO_*` settings. If all dedicated connection fields are blank, they may use `GOREST_DB_*` only when `GOREST_APP_ENV` is development, test, or testing and the MongoDB driver is enabled. In either case, tests refuse production mode and require `GOREST_TEST_MONGO_CONFIRM_DATABASE` to exactly match the target database. Names must end in `_test` or `_testing`, unless `GOREST_TEST_MONGO_ALLOW_NON_TEST_DATABASE` repeats the exact approved non-production name.

The test suite applies the up migration, which creates the collection and indexes if needed and records migration state, and runs the sample seeder once. It does not drop the database or run the down migration. Three sample documents with placeholder password hashes and migration/seeder history remain in the database; tests insert unique user documents and clean up only their own IDs/emails. Use a disposable, dedicated non-production database and an account limited to the required migration-metadata, collection, index, and document permissions. See the [testing guide](./testing.md).
