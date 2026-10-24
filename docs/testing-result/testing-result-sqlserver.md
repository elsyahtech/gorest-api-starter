# SQL Server test result

**Result:** `PASS`  
**Run date:** October 23, 2026  
**Target:** local container, `127.0.0.1:1433`, disposable test database; isolated disposable target cleaned up after the run  
**Server:** SQL Server 2022 (RTM-CU26), 16.0.4265.3 (X64)  
**Go / driver:** `go1.26.5 darwin/amd64`; `github.com/microsoft/go-mssqldb v1.11.2`  
**Source commit:** `21e0662bc6d4eb24ffdaf68d9ef57e158ade8523` (tests rerun from a clean archive of this commit)
**Command:** `GOREST_DRIVER_INTEGRATION=sqlserver go test -v -count=1 -timeout=20m ./test/drivers -run '^TestDatabaseDriverIntegration$'`

Credentials are supplied through environment variables and are not recorded here. The disposable test target was cleaned up after this run.

| Scenario | Result and evidence |
| :------- | :------------------ |
| Migration and seeder | `passed` — migrations and seeders completed; 11 seeded users were read. |
| ORM CRUD | `passed` — create, read, update, and delete completed. |
| Transaction rollback | `passed` — normal rollback passed; a constraint failure on the second write rolled back the earlier write. |
| Mid-step rollback | `passed` — a constraint failure on the second transaction write rolled back the earlier write. |
| Conflicts and invalid operations | `passed` — duplicate request key, malformed JSON, and update of a missing user were rejected without creating rows. |
| Concurrent requests | `passed` — all eight simultaneous creates returned HTTP `200`, `success=true`, and matching persisted rows. |
| Startup failure | `passed` — the driver integration run rejected an invalid migration; a subprocess run against a reserved then released local port exited with `failed initializing database`. |
| Shutdown | `passed` — per-driver `app.Close()` closed the database and logger; the app-level E2E shutdown test waited for an in-flight request to complete before exit. |

Recorded output:

```text
migration and seeder passed; read 11 seeded users
sql server ORM create/read/update/delete passed
sql server transaction rollback passed
sql server mid-step constraint failure rolled back the earlier write
duplicate request key, malformed JSON, and missing-user update were rejected without creating rows
sql server concurrent HTTP creates returned HTTP 200, success=true, and matching persisted rows for all 8 requests
sql server startup rejected an invalid migration
sql server application shutdown passed
PASS
```

Supplemental lifecycle evidence ran October 23, 2026 against the working tree based on the listed source commit: `GOREST_DRIVER_LIFECYCLE_DRIVER=sqlserver go test -v -count=1 -timeout=2m ./test/drivers -run '^TestDatabaseDriverStartupFailsWithoutDatabase$'` passed. The shared app-level lifecycle E2E selection also passed as recorded in [the result index](../testing-result.md). These supplemental test-harness changes are not yet committed.

Supplemental full target evidence ran October 24, 2026 with `make test-e2e-sqlserver`. The disposable SQL Server 2022 container passed driver integration and the complete shared route/lifecycle E2E suite, including native SQL CRUD, preload, and pagination. This run exposed SQL Server's rejection of `FETCH NEXT 0`; the model now returns an empty result directly for a zero limit, and the scenario passed. The target removed its container and network after completion.

This evidence applies to the listed server version, driver, source revision, and test configuration only.
