# MongoDB test result

**Result:** `PASS`  
**Run date:** October 23, 2026  
**Target:** local container, `127.0.0.1:27017`, disposable test database; isolated disposable target cleaned up after the run  
**Server:** 8.2.12  
**Go / driver:** `go1.26.5 darwin/amd64`; `go.mongodb.org/mongo-driver/v2 v2.9.2`  
**Source:** `1d86605f56cff6c5df751aabdb93bddc995b4870`; working tree contained uncommitted changes  
**Command:** `GOREST_DRIVER_INTEGRATION=mongo go test -v -count=1 -timeout=20m ./test/drivers -run '^TestDatabaseDriverIntegration$'`

Credentials are supplied through environment variables and are not recorded here. The disposable test target was cleaned up after this run.

| Scenario | Result and evidence |
| :------- | :------------------ |
| Migration and seeder | `passed` — migrations and seeders completed; 11 seeded users were read. |
| ORM CRUD | `passed` — create, read, update, and delete completed. |
| Transaction rollback | `not applicable` — transaction atomicity is not provided by the GOREST ORM transaction API exercised here. |
| Mid-step rollback | `not applicable` — the tested ORM API does not provide this driver transaction atomicity. |
| Conflicts and invalid operations | `passed` — duplicate request key, malformed JSON, and update of a missing user were rejected without creating rows. |
| Concurrent requests | `passed` — all eight simultaneous creates returned HTTP `200`, `success=true`, and matching persisted rows. |
| Startup failure | `partial` — startup rejected an invalid migration after connecting. Initial connection failure was not exercised for this driver. |
| Shutdown | `partial` — normal app shutdown passed; interrupted requests and shutdown timeout were not exercised by this driver test. |

Recorded output:

```text
migration and seeder passed; read 11 seeded users
mongodb ORM create/read/update/delete passed
mongodb transaction atomicity: not applicable to the starter's GOREST ORM transaction API
duplicate request key, malformed JSON, and missing-user update were rejected without creating rows
mongodb concurrent HTTP creates returned HTTP 200, success=true, and matching persisted rows for all 8 requests
mongodb startup rejected an invalid migration
mongodb application shutdown passed
PASS
```

This evidence applies to the listed server version, driver, source revision, and test configuration only.
