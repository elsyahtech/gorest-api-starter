# PostgreSQL test result

**Result:** `PASS`  
**Run date:** October 23, 2026  
**Target:** local PostgreSQL container, `127.0.0.1:5432`; isolated disposable database created for the run and removed afterward  
**Server:** PostgreSQL 16.3 (`x86_64-pc-linux-musl`)  
**Go / driver:** `go1.26.5 darwin/amd64`; `github.com/lib/pq v1.12.3`  
**Source commit:** `21e0662bc6d4eb24ffdaf68d9ef57e158ade8523` (tests rerun from a clean archive of this commit)
**Command:** `GOREST_DRIVER_INTEGRATION=postgres go test -v -count=1 -timeout=20m ./test/drivers -run '^TestDatabaseDriverIntegration$'`

Credentials are supplied through environment variables and are not recorded here. The test target was disposable and cleaned up after this run.

| Scenario | Result and evidence |
| :------- | :------------------ |
| Migration and seeder | `passed` — migrations and seeders completed; 11 seeded users were read. |
| ORM CRUD | `passed` — create, read, update, and delete completed. |
| Transaction rollback | `passed` — transaction rollback completed; a constraint failure on the second write rolled back the earlier write. |
| Conflicts and invalid operations | `passed` — duplicate request key, malformed JSON, and update of a missing user were rejected without creating rows. |
| Concurrent requests | `passed` — all eight simultaneous creates returned HTTP `200`, `success=true`, and matching persisted rows. |
| Startup failure | `passed` — the driver integration run rejected an invalid migration; a subprocess run against a reserved then released local port exited with `failed initializing database`. |
| Shutdown | `passed` — per-driver `app.Close()` closed the database and logger; the app-level E2E shutdown test waited for an in-flight request to complete before exit. |

Recorded output:

```text
migration and seeder passed; read 11 seeded users
postgres ORM create/read/update/delete passed
postgres transaction rollback passed
postgres mid-step constraint failure rolled back the earlier write
duplicate request key, malformed JSON, and missing-user update were rejected without creating rows
postgres concurrent HTTP creates returned HTTP 200, success=true, and matching persisted rows for all 8 requests
postgres startup rejected an invalid migration
postgres application shutdown passed
PASS
```

Supplemental lifecycle evidence ran October 23, 2026 against the working tree based on the listed source commit: `GOREST_DRIVER_LIFECYCLE_DRIVER=postgres go test -v -count=1 -timeout=2m ./test/drivers -run '^TestDatabaseDriverStartupFailsWithoutDatabase$'` passed. The shared app-level lifecycle E2E selection also passed as recorded in [the result index](../testing-result.md). These supplemental test-harness changes are not yet committed.

This evidence applies to the listed PostgreSQL version, driver, source revision, and test configuration only.

## Full route and lifecycle suite

**Run date:** October 24, 2026

**Target:** disposable PostgreSQL 16.3 container, isolated Compose project and network

**Command:** `make test-e2e-postgres`

**Result:** `PASS` on the current working tree.

The target passed the complete in-process route and lifecycle suite against PostgreSQL, including department and order routes, ORM user queries and writes, transactions, bulk operations, native SQL create/read/update/delete, database-startup failure, migration failure, and graceful shutdown. MongoDB and ScyllaDB native route cases are excluded from this SQL-driver run. The runner removed its PostgreSQL container and Compose network after the run. This supplemental evidence is from the working tree and is not part of the source commit listed above.
