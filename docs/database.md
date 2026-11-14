# SQLite configuration

This starter uses SQLite through GOREST's database service. Runtime settings are read by `app/config/database.go`.

## Application database

The default settings are:

| Variable | Default | Purpose |
|---|---|---|
| `GOREST_DB_DRIVER` | `sqlite` | The SQLite-only driver selected by this starter. |
| `GOREST_DB_NAME` | `gorest` | Database file name, without the `.db` extension. |
| `GOREST_DB_DIRECTORY` | `./app/database` | Directory containing the database file and schema assets. |
| `GOREST_DB_MIGRATION` | `false` | Enables migrations during application startup. |
| `GOREST_DB_SEEDER` | `false` | Enables sample data insertion during startup. |
| `GOREST_DB_MIGRATIONS_DIRECTORY` | `migrations` | Migration asset directory, resolved under `GOREST_DB_DIRECTORY`. |
| `GOREST_DB_SEEDERS_DIRECTORY` | `seeders` | Seeder asset directory, resolved under `GOREST_DB_DIRECTORY`. |

The resulting default database file is `app/database/gorest.db`. The starter rejects any driver other than `sqlite`. It does not require network credentials or a separately provisioned database server.

The `users` table is created by `app/database/migrations/sqlite/001_create_users_table.sql`. Set `GOREST_DB_MIGRATION=true` for a controlled migration run, then turn it off for normal application startup. The down migration drops the `users` table; never run it against data that must be retained. The sample seeder inserts ten rows with placeholder password hashes; these are not usable login credentials. Enable it only when sample data is useful.

The database directory must be writable by the application. Back up the database file using a method appropriate for SQLite, and avoid copying a live database file without considering SQLite's journal/WAL state.

## Integration-test database

SQLite integration tests create a unique database in `testing.T.TempDir()` for each test. The tests point migration/seeder lookup to the repository's SQLite assets and do not read `.env`, connect to a network database, or write to the default `app/database/gorest.db`.

The HTTP integration test enables migration and seeding and verifies the user API against its temporary database. The suite does not truncate or drop any application database. Optional Redis integration tests are separate; see the [testing guide](./testing.md).

See the [GOREST configuration guide](https://github.com/elsyahtech/gorest/blob/v1.1.21/docs/configuration.md) and [testing guide](./testing.md) for runtime and test configuration.
