# MongoDB test result

**Result:** `PASS`  
**Run date:** October 23, 2026  
**Target:** local container, `127.0.0.1:27017`, disposable test database; isolated disposable target cleaned up after the run  
**Server:** 8.2.12  
**Go / driver:** `go1.26.5 darwin/amd64`; `go.mongodb.org/mongo-driver/v2 v2.9.2`  
**Source commit:** `21e0662bc6d4eb24ffdaf68d9ef57e158ade8523` (tests rerun from a clean archive of this commit)
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
| Startup failure | `passed` — the driver integration run rejected an invalid migration; a subprocess run against a reserved then released local port exited with `failed initializing database`. |
| Shutdown | `passed` — per-driver `app.Close()` closed the database and logger; the app-level E2E shutdown test waited for an in-flight request to complete before exit. |

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

Supplemental lifecycle evidence ran October 23, 2026 against the working tree based on the listed source commit: `GOREST_DRIVER_LIFECYCLE_DRIVER=mongo go test -v -count=1 -timeout=2m ./test/drivers -run '^TestDatabaseDriverStartupFailsWithoutDatabase$'` passed. The shared app-level lifecycle E2E selection also passed as recorded in [the result index](../testing-result.md). These supplemental test-harness changes are not yet committed.

Supplemental E2E evidence ran October 24, 2026 with `make test-e2e-mongo`. The disposable MongoDB 8.2.12 container passed migration/seeder, ORM CRUD, conflicts, concurrency, startup failure, and shutdown checks. The shared ORM route/lifecycle suite passed for applicable scenarios, including user, department, and order routes; transaction-atomicity cases were omitted because the exercised GOREST ORM API does not provide that behavior. All four native MongoDB route scenarios (create, read, update, delete) passed in their separate stage. The target removed its container and network after completion.

MongoDB's suite uses MongoDB-specific fixtures and collection cleanup. Its shared ORM route coverage is distinct from the four native MongoDB query routes, and both stages ran successfully in the recorded target.

This evidence applies to the listed server version, driver, source revision, and test configuration only.
