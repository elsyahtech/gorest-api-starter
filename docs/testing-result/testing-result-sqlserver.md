# SQL Server test result

**Result:** `PASS`  
**Run date:** October 23, 2026  
**Target:** local container, `127.0.0.1:1433`, disposable test database; isolated disposable target cleaned up after the run  
**Server:** SQL Server 2022 (RTM-CU26), 16.0.4265.3 (X64)  
**Go / driver:** `go1.26.5 darwin/amd64`; `github.com/microsoft/go-mssqldb v1.11.2`  
**Source:** `1d86605f56cff6c5df751aabdb93bddc995b4870`; working tree contained uncommitted changes  
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
| Startup failure | `partial` — startup rejected an invalid migration after connecting. Initial connection failure was not exercised for this driver. |
| Shutdown | `partial` — normal app shutdown passed; interrupted requests and shutdown timeout were not exercised by this driver test. |

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

This evidence applies to the listed server version, driver, source revision, and test configuration only.
