# SQLite test result

**Result:** `PASS`  
**Run date:** October 23, 2026  
**Target:** temporary SQLite database managed by Go test cleanup  
**Command:** `GOREST_DRIVER_INTEGRATION=sqlite go test ./test/drivers -run '^TestDatabaseDriverIntegration$' -count=1 -v`
**SQLite engine:** `3.53.4` (`modernc.org/sqlite v1.60.1`)  
**Go:** `go1.26.5 darwin/amd64`  
**Source commit:** `21e0662bc6d4eb24ffdaf68d9ef57e158ade8523` (tests rerun from a clean archive of this commit)

| Scenario | Result and evidence |
| :------- | :------------------ |
| Migration and seeder | `passed` — migrations and seeders completed; 11 seeded users were read. |
| ORM CRUD | `passed` — create, read, update, and delete completed. |
| Transaction rollback | `passed` — transaction rollback completed. A constraint failure on the second write rolled back the earlier write. |
| Conflicts and invalid operations | `passed` — duplicate request key, malformed JSON, and update of a missing user were rejected without creating rows. |
| Concurrent requests | `passed` for the bounded scenario — all 8 simultaneous creates returned HTTP `200`, `success=true`, and matching persisted rows. The SQLite test pool uses one connection to queue writes. |
| Startup failure | `passed` — startup tasks rejected an invalid migration after connecting to the temporary SQLite database. A network connection failure is not applicable to embedded SQLite. |
| Shutdown | `passed` — per-driver `app.Close()` closed the database and logger; the app-level E2E shutdown test waited for an in-flight request to complete before exit. |

Recorded output:

```text
migration and seeder passed; read 11 seeded users
sqlite ORM create/read/update/delete passed
sqlite transaction rollback passed
sqlite mid-step constraint failure rolled back the earlier write
duplicate request key, malformed JSON, and missing-user update were rejected without creating rows
sqlite concurrent HTTP creates returned HTTP 200, success=true, and matching persisted rows for all 8 requests
sqlite startup rejected an invalid migration
sqlite application shutdown passed
PASS
```

The SQLite integration rerun, including invalid-migration rejection and normal close, passed on October 23, 2026 from the working tree based on the listed source commit. The shared app-level lifecycle E2E selection also passed as recorded in [the result index](../testing-result.md). These supplemental test-harness changes are not yet committed. The default driver package skips native database integration unless `GOREST_DRIVER_INTEGRATION` is set. These results apply to this source revision and its temporary SQLite configuration.
