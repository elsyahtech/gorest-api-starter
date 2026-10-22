# Database, migrations, and seeders

The application database settings are in `app/config/database.go`. The database is enabled and PostgreSQL is selected by default. The credentials in that file are example values included in source code. Replace them for local use, and do not store production credentials in the repository.

## Drivers and file locations

The configuration lists these drivers: MySQL, PostgreSQL, SQL Server, Oracle, SQLite, MongoDB, and ScyllaDB. Migration and seeder files are grouped by driver:

```text
app/database/
  migrations/<driver>/
  seeders/<driver>/
```

Migrations define the `departments`, `users`, and `orders` tables or collections. Seeders provide sample data for those resources. The PostgreSQL examples include departments with IDs `dept-01` through `dept-05`, sample users, and sample orders. Check the files for your selected driver before relying on specific IDs or data.

## Paths and startup settings

`app/config/database.go` includes settings for:

- `Driver`, host, port, username, password, and database name.
- TLS and certificate or wallet paths, when needed.
- `DatabaseDirectory`, `MigrationDirectory`, and `SeederDirectory`.
- Timeouts and connection pool limits.
- `Migration` and `Seeder`. Both are enabled by default, so their files run when the application starts.

Directory paths are relative to the process working directory. Run the application from the repository root as described in [Getting started](./getting-started.md), or adjust the paths if your deployment layout is different. In production, disable automatic migrations and seeders on application replicas. Run schema and data changes as a controlled deployment step instead.

GOREST runs migrations before seeders and does not start the application if a startup migration or seeder fails. A failed script may have already changed the database, so inspect the database and migration history before retrying it. Treat applied migration and seeder files as immutable; add a new file for later changes instead of editing an applied file. For the framework's one-time database startup task workflow and recovery details, see the [GOREST v1.0.2 migration and seeder guide](https://github.com/elsyahtech/gorest/blob/v1.0.2/docs/guide.md).

## Example relationships

- Users store `department_id`; the SQL schema also creates an index for this column.
- Orders store `user_id`; the PostgreSQL schema defines a foreign key from orders to users.
- The PostgreSQL SQL schema requires user email addresses to be unique.

Check the migration files for your selected driver before assuming that constraints or behavior are the same across all databases.

## Native queries and transactions

The native SQL user routes use PostgreSQL syntax and `$1` placeholders. Adapt the query and placeholders before using a different SQL driver. Use the MongoDB routes only with MongoDB, and the Scylla routes only with ScyllaDB.

Some ORM operations use `BeginTx`, `Commit`, and `Rollback`. In these examples, GOREST transaction calls are no-ops for MongoDB and ScyllaDB. Do not assume that writes using those drivers are atomic transactions.

For native route details, see the [API reference](./api-reference.md). For driver settings and ORM APIs, see the [GOREST v1.0.2 documentation](https://github.com/elsyahtech/gorest/tree/v1.0.2/docs).
