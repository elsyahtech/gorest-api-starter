# Oracle test result

**Result:** `PASS` — all three repeat runs passed the strict concurrency assertion after aligning the Fiber test-client timeout with the database operation timeout
**Run date:** October 23, 2026  
**Target:** local container, `127.0.0.1:1521/FREEPDB1`, temporary `GOREST_IT_<run>` schema; dropped after the run  
**Server:** Oracle AI Database 26ai Free Release 23.26.3.0.0  
**Go / driver:** `go1.26.5 darwin/amd64`; `github.com/godror/godror v0.51.5`  
**Source:** application baseline commit `21e0662bc6d4eb24ffdaf68d9ef57e158ade8523`; the uncommitted test harness changes increase the Fiber test-client timeout from 10 to 30 seconds and add the startup-failure check
**Command:** `GOREST_DRIVER_INTEGRATION=oracle go test -v -count=1 -timeout=20m ./test/drivers -run '^TestDatabaseDriverIntegration$'`

Credentials are supplied through environment variables and are not recorded here. `GOREST_DRIVER_LIB_DIR` pointed to the host's Oracle Instant Client 19.16 directory. The temporary schema used test-only grants and was dropped after this run.

| Scenario | Result and evidence |
| :------- | :------------------ |
| Migration and seeder | `passed` — migrations and seeders completed; 11 seeded users were read. |
| ORM CRUD | `passed` — create, read, update, and delete completed. |
| Transaction rollback | `passed` — normal rollback passed; a constraint failure on the second write rolled back the earlier write. |
| Mid-step rollback | `passed` — a constraint failure on the second transaction write rolled back the earlier write. |
| Conflicts and invalid operations | `passed` — duplicate request key, malformed JSON, and update of a missing user were rejected without creating rows. |
| Concurrent requests | `passed` — three consecutive runs each returned HTTP `200`, `success=true`, and matching persisted rows for all eight simultaneous creates. The Fiber test-client timeout is 30 seconds, matching the database operation timeout. |
| Startup failure | `passed` — the driver integration run rejected an invalid migration; a subprocess run against a reserved then released local port exited with `failed initializing database`. |
| Shutdown | `passed` — per-driver `app.Close()` closed the database and logger; the app-level E2E shutdown test waited for an in-flight request to complete before exit. |

The earlier sequential run recorded below used the 10-second Fiber test-client timeout and failed four requests without a response. The test client deadline was shorter than the configured 30-second database operation timeout, so the test harness timeout was aligned to 30 seconds. The strict response and persisted-row assertions were retained:

```text
POST user c-oracle-18e116b815e92310-01 failed without a response and no matching row was found: test: got empty response
POST user c-oracle-18e116b815e92310-05 failed without a response and no matching row was found: test: got empty response
POST user c-oracle-18e116b815e92310-04 failed without a response and no matching row was found: test: got empty response
POST user c-oracle-18e116b815e92310-03 failed without a response and no matching row was found: test: got empty response
only 4 of 8 valid concurrent creates returned a consistent success response
FAIL
```

After that change, three consecutive runs passed. Each run also passed migration/seeder, ORM CRUD, transaction and mid-step rollback, conflict, and invalid-migration startup checks. Runs used temporary Oracle schemas, which were dropped after each test.

Supplemental Oracle startup-failure evidence ran October 23, 2026 against the working tree based on the listed source commit: `GOREST_DRIVER_LIFECYCLE_DRIVER=oracle GOREST_DRIVER_LIB_DIR=/Users/elmansyah/oracle/instantclient_19_16 go test -v -count=1 -timeout=2m ./test/drivers -run '^TestDatabaseDriverStartupFailsWithoutDatabase$'` passed. The shared app-level lifecycle E2E selection also passed as recorded in [the result index](../testing-result.md). These supplemental test-harness changes are not yet committed.

Supplemental full target evidence ran October 24, 2026 with `make test-e2e-oracle`. Oracle Free 26ai in the disposable Compose container passed the driver integration scenarios and the complete shared route/lifecycle E2E suite, including native SQL CRUD. The run also passed the preload-with-missing-department and zero-limit query cases. The target downloaded and used its temporary Oracle Instant Client package, then removed the container, network, and temporary client files on exit.

This evidence applies to the listed server version, driver, source revision, and test configuration only.
