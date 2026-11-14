# Oracle configuration

This starter uses Oracle through GOREST's database service. Runtime settings are read by `app/config/database.go`.

## Application connection

Set these variables in `.env`:

| Variable | Purpose |
|---|---|
| `GOREST_DB_DRIVER` | Must be `oracle` for this starter. |
| `GOREST_DB_HOST`, `GOREST_DB_PORT` | Oracle listener host and TCP port (default `1521`). |
| `GOREST_DB_NAME` | Oracle service name, not a SID; required. |
| `GOREST_DB_USER`, `GOREST_DB_PASSWORD` | Application schema username and password; both required. |
| `GOREST_DB_LIB_DIR_PATH` | Optional Oracle Instant Client library directory when it is not available in the system library path. |
| `GOREST_DB_TLS` | `disable` for TCP or `enable` for TCPS. Remote production connections should use TCPS. |
| `GOREST_DB_CA_PATH` | Oracle wallet directory (`TNS_ADMIN`); required when `GOREST_DB_TLS=enable`. This is not a PEM CA file. |
| `GOREST_DB_MIGRATION` | Enables Oracle migrations during application startup. |
| `GOREST_DB_SEEDER` | Enables sample data insertion; disabled by default. |

Install a compatible Oracle Instant Client on the application host. Set `GOREST_DB_LIB_DIR_PATH` to its library directory if the client is not discoverable through the operating system's library path. The starter accepts only the Oracle driver and requires an explicit service name and credentials.

For encrypted Oracle Net connections, set `GOREST_DB_TLS=enable` and point `GOREST_DB_CA_PATH` to a wallet directory containing the Oracle Net configuration and wallet files. GOREST uses this path as the Oracle client configuration directory (`TNS_ADMIN`). Do not supply generic PEM paths through `GOREST_DB_CERT_PATH` or `GOREST_DB_KEY_PATH`; Oracle uses its wallet configuration. Keep server certificate verification enabled in the wallet and Oracle Net configuration. `GOREST_DB_TLS=disable` uses unencrypted TCP and is intended only for trusted local or isolated testing environments.

The `users` table is created by `app/database/migrations/oracle/001_create_users_table.sql`; the active-user index is created by the next migration. Set `GOREST_DB_MIGRATION=true` for a controlled migration run, then turn it off for normal application startup. The down migrations remove the index and table; never run them against data that must be retained. The sample seeder inserts ten rows with placeholder password hashes; these are not usable login credentials. Enable it only when sample data is useful.

Use a dedicated application schema with only the required privileges. Applying migrations and seeders requires permission to create the GOREST migration lock and migration/seeder tracking objects, create the users table and index, and insert sample users. The schema also needs quota on its target tablespace. Do not run these operations as `SYS`, `SYSTEM`, or another administrative account.

## Integration-test schema

The Oracle integration suite reads `GOREST_DB_*` settings from the process environment or `.env`. It requires `GOREST_TEST_ORACLE_CONFIRM_USER` to exactly match `GOREST_DB_USER` and refuses production mode and built-in administrative accounts. The test schema needs `CREATE TABLE` and tablespace quota so it can create the framework lock/history objects and starter tables. The HTTP flow runs migrations and seeders, then inserts and removes only its UUID-scoped test user. The library test also cleans up its unique test user. The suite does not truncate tables or drop the schema, but migrations, migration history, and ten seeded sample rows persist. Use a disposable, dedicated non-production schema.

See the [GOREST Oracle guide](https://github.com/elsyahtech/gorest/blob/v1.1.21/docs/testing-result/testing-result-oracle.md) and [testing guide](./testing.md) for client, wallet, and integration-test setup.
