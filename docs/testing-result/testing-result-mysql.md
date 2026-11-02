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

## Full route and lifecycle suite

**Run date:** October 24, 2026

**Target:** disposable MySQL 26.7.0 container, isolated Compose project and network

**Command:** `make test-e2e-mysql`

**Result:** `PASS` on the current working tree.

The target passed the complete in-process route and lifecycle suite against MySQL, including department and order routes, ORM user queries and writes, transactions, bulk operations, native SQL create/read/update/delete, database-startup failure, migration failure, and graceful shutdown. MongoDB and ScyllaDB native route cases are excluded from this SQL-driver run. The runner removed its MySQL container and Compose network after the run. This supplemental evidence is from the working tree and is not part of the source commit listed above.

## Local MySQL run including login E2E

**Run date:** November 1, 2026

**Target:** local MySQL container at `127.0.0.1:3306`, using a dedicated disposable schema that was removed after testing

**Result:** `PASS` for MySQL database integration and the MySQL E2E suite.

The login scenario is now driver-independent and was run against both SQLite and MySQL. It verifies that an unauthenticated request is rejected, an incorrect password is rejected, valid credentials return an access cookie, and that cookie grants access to the protected user route.

The complete MySQL E2E suite passed, as did `TestDatabaseDriverIntegration`, application/root unit tests, `go vet ./...`, and the repository baseline linters for the database/E2E test packages. The MySQL E2E run excluded native CQL and Mongo query cases; Redis E2E remained skipped because `GOREST_REDIS_E2E=1` was not enabled. These are separate services/features, not MySQL scenarios. This evidence shows the listed MySQL flows passed on this local setup; it does not claim exhaustive coverage of every application branch or production environment.

## External MySQL run with verified TLS and mTLS

**Run date:** November 1, 2026

**Target:** the external MySQL database provided for testing; credentials and certificate contents are intentionally not recorded here

**TLS:** a test-only `go-sql-driver/mysql` profile verified the server certificate chain and hostname, and presented the configured client certificate for mTLS. TLS certificate verification remained enabled.

**Result:** `PASS` for `TestDatabaseDriverIntegration` and the complete MySQL E2E suite, including login, protected-route access, CRUD, native SQL, transaction rollback, request validation, concurrency, startup failure, and lifecycle scenarios. The ordinary `go test ./... -count=1` suite also passed with the default local setup.

One initial full E2E attempt encountered a remote connection reset late in the suite. The failing E2E scenario passed when rerun individually, and the complete suite passed on a subsequent full run. MongoDB, ScyllaDB, and Redis-specific E2E scenarios remained skipped because this run was scoped to MySQL. The test target is disposable and the suite ran migrations, seeders, and data-changing scenarios against it.

## External MySQL storage endpoint and 500-row measurement

**Run date:** November 1, 2026

**Target:** the external MySQL test database, connected with the same verified TLS/mTLS test profile described above.

**Command:** `GOREST_DRIVER_INTEGRATION=mysql go test -v -count=1 -timeout=5m ./test/drivers -run '^TestMySQLStorageEndpointIntegration$'`

**Result:** `PASS`. The test registered `/gorest/storage` on an isolated test server, inserted 500 rows with 4,000-byte payloads into a uniquely named temporary InnoDB table, and removed that table during cleanup.

| Measurement | `used_bytes` |
| :---------- | -----------: |
| Before inserting the 500 rows | 196,608 |
| After inserting the 500 rows | 245,760 |
| Observed increase | 49,152 |

These are MySQL `information_schema` table data/index-size estimates, not an exact accounting of payload bytes; the increase should not be interpreted as the precise physical size of the inserted data. The test-only route registration does not enable monitoring endpoints in the normal application configuration.

A separate attempt to rerun the full database integration test against this external target stopped at its existing seeder assertion (`1` user found; at least `10` expected), before reaching the other integration scenarios. The isolated storage endpoint test above does not depend on seed data and passed.

## Latest disposable MySQL rerun

**Run date:** November 2, 2026

**Target:** temporary MySQL `26.7.0` container and isolated Compose network, automatically removed after the run.

**Command:** `make test-e2e-mysql`

**Result:** `PASS`. The database integration tests completed migrations and seeders (11 users read), ORM and native SQL CRUD, transaction and mid-step rollback, invalid/conflicting requests, eight concurrent creates, startup and migration failure checks, and shutdown. The complete MySQL E2E suite also passed, including login and protected-route access, departments, orders, and user routes. This run used a fresh disposable local database; it does not change the separate external MySQL test result above.
