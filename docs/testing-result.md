# Database driver test results

The evidence pages record the automated driver integration runs performed on October 23, 2026. The full database scenarios were run against disposable local databases/containers with `go test -v -count=1 -timeout=20m ./test/drivers -run '^TestDatabaseDriverIntegration$'`; Docker Compose was not used. Supplemental startup-failure checks ran separately for each network driver against an unavailable local port, and lifecycle E2E checks verified migration failure and graceful shutdown. Credentials are omitted and temporary targets were cleaned up. Each page records the server version, Go toolchain, driver module version, source provenance, scenarios, and observed output.

| Driver | Result |
| :----- | :----- |
| PostgreSQL | [Test Result](./testing-result/testing-result-postgres.md) |
| MySQL | [Test Result](./testing-result/testing-result-mysql.md) |
| Oracle | [Test Result](./testing-result/testing-result-oracle.md) |
| SQL Server | [Test Result](./testing-result/testing-result-sqlserver.md) |
| SQLite | [Test Result](./testing-result/testing-result-sqlite.md) |
| MongoDB | [Test Result](./testing-result/testing-result-mongo.md) |
| ScyllaDB | [Test Result](./testing-result/testing-result-scylla.md) |

The [driver verification matrix](./testing.md#driver-verification-matrix) summarizes the scenarios and their status. A clean-archive sweep of commit `21e0662bc6d4eb24ffdaf68d9ef57e158ade8523` passed the strict concurrency contract for all seven drivers. Oracle then passed three repeat runs after the Fiber test-client timeout was aligned to the database operation timeout; its evidence page records that test-harness source adjustment and the earlier timeout failure. MongoDB and ScyllaDB transaction rollback is not applicable to the GOREST ORM transaction API exercised here. The SQLite result page records the bundled SQLite engine and Go driver versions.

The full driver baseline used a clean archive of commit `21e0662bc6d4eb24ffdaf68d9ef57e158ade8523`. The supplemental Oracle timeout and lifecycle test-harness changes were run from the working tree and are not yet part of that commit; their provenance is stated on the relevant evidence pages. The lifecycle E2E selection passed on October 23 with `go test -v -count=1 -timeout=3m ./test/e2e -run '^Test(DatabaseConnectionFailureAbortsStartup|MigrationFailureAbortsStartupTasks|ShutdownWaitsForInFlightRequest)$'`. These results describe the specific local server versions and configurations recorded; they do not establish compatibility with other versions, topologies, or deployment configurations.
