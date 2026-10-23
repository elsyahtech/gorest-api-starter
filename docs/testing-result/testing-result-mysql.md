# MySQL test result

**Result:** `PASS`  
**Run date:** October 23, 2026  
**Target:** local MySQL container, `127.0.0.1:3306`; isolated disposable database created for the run and removed afterward  
**Server:** MySQL Community Server `26.7.0`  
**Go / driver:** `go1.26.5 darwin/amd64`; `github.com/go-sql-driver/mysql v1.10.1`  
**Source:** `1d86605f56cff6c5df751aabdb93bddc995b4870`; working tree contained uncommitted changes  
**Command:** `GOREST_DRIVER_INTEGRATION=mysql go test -v -count=1 -timeout=20m ./test/drivers -run '^TestDatabaseDriverIntegration$'`

Credentials are supplied through environment variables and are not recorded here. The test target was disposable and cleaned up after this run.

| Scenario | Result and evidence |
| :------- | :------------------ |
| Migration and seeder | `passed` — migrations and seeders completed; 11 seeded users were read. |
| ORM CRUD | `passed` — create, read, update, and delete completed. |
| Transaction rollback | `passed` — transaction rollback completed; a constraint failure on the second write rolled back the earlier write. |
| Conflicts and invalid operations | `passed` — duplicate request key, malformed JSON, and update of a missing user were rejected without creating rows. |
| Concurrent requests | `passed` — all eight simultaneous creates returned HTTP `200`, `success=true`, and matching persisted rows. |
| Startup failure | `partial` — startup rejected an invalid migration after connecting. Initial connection failure was not exercised for this driver. |
| Shutdown | `partial` — normal app shutdown passed; interrupted requests and shutdown timeout were not exercised by this driver test. |

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

This evidence applies to the listed MySQL version, driver, source revision, and test configuration only.
