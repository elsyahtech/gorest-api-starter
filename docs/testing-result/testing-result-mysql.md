# MySQL test result

**Result:** `PASS`  
**Run date:** October 23, 2026  
**Target:** local MySQL container, `127.0.0.1:3306`; isolated disposable database created for the run and removed afterward  
**Server:** MySQL Community Server `26.7.0`  
**Go / driver:** `go1.26.5 darwin/amd64`; `github.com/go-sql-driver/mysql v1.10.1`  
**Source commit:** `21e0662bc6d4eb24ffdaf68d9ef57e158ade8523` (tests rerun from a clean archive of this commit)
**Command:** `GOREST_DRIVER_INTEGRATION=mysql go test -v -count=1 -timeout=20m ./test/drivers -run '^TestDatabaseDriverIntegration$'`

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
mysql ORM create/read/update/delete passed
mysql transaction rollback passed
mysql mid-step constraint failure rolled back the earlier write
duplicate request key, malformed JSON, and missing-user update were rejected without creating rows
mysql concurrent HTTP creates returned HTTP 200, success=true, and matching persisted rows for all 8 requests
mysql startup rejected an invalid migration
mysql application shutdown passed
PASS
```

Supplemental lifecycle evidence ran October 23, 2026 against the working tree based on the listed source commit: `GOREST_DRIVER_LIFECYCLE_DRIVER=mysql go test -v -count=1 -timeout=2m ./test/drivers -run '^TestDatabaseDriverStartupFailsWithoutDatabase$'` passed. The shared app-level lifecycle E2E selection also passed as recorded in [the result index](../testing-result.md). These supplemental test-harness changes are not yet committed.

This evidence applies to the listed MySQL version, driver, source revision, and test configuration only.
