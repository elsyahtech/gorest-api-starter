# Oracle test result

**Result:** `PASS`  
**Run date:** October 23, 2026  
**Target:** local container, `127.0.0.1:1521/FREEPDB1`, temporary `GOREST_IT_<run>` schema; dropped after the run  
**Server:** Oracle AI Database 26ai Free Release 23.26.3.0.0  
**Go / driver:** `go1.26.5 darwin/amd64`; `github.com/godror/godror v0.51.5`  
**Source:** `1d86605f56cff6c5df751aabdb93bddc995b4870`; working tree contained uncommitted changes  
**Command:** `GOREST_DRIVER_INTEGRATION=oracle go test -v -count=1 -timeout=20m ./test/drivers -run '^TestDatabaseDriverIntegration$'`

Credentials are supplied through environment variables and are not recorded here. `GOREST_DRIVER_LIB_DIR` pointed to the host's Oracle Instant Client 19.16 directory. The temporary schema used test-only grants and was dropped after this run.

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
oracle ORM create/read/update/delete passed
oracle transaction rollback passed
oracle mid-step constraint failure rolled back the earlier write
duplicate request key, malformed JSON, and missing-user update were rejected without creating rows
oracle concurrent HTTP creates returned HTTP 200, success=true, and matching persisted rows for all 8 requests
oracle startup rejected an invalid migration
oracle application shutdown passed
PASS
```

This evidence applies to the listed server version, driver, source revision, and test configuration only.
