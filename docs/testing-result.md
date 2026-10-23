# Database driver test results

The following evidence pages record the automated driver integration runs performed on October 23, 2026. Each driver was tested against a disposable local database/container with `go test -v -count=1 -timeout=20m ./test/drivers -run '^TestDatabaseDriverIntegration$'`; Docker Compose was not used. Credentials are omitted and the temporary targets were cleaned up. Each page records the server version, Go toolchain, driver module version, source commit, scenarios, and observed output.

| Driver | Result |
| :----- | :----- |
| PostgreSQL | [Test Result](./testing-result/testing-result-postgres.md) |
| MySQL | [Test Result](./testing-result/testing-result-mysql.md) |
| Oracle | [Test Result](./testing-result/testing-result-oracle.md) |
| SQL Server | [Test Result](./testing-result/testing-result-sqlserver.md) |
| SQLite | [Test Result](./testing-result/testing-result-sqlite.md) |
| MongoDB | [Test Result](./testing-result/testing-result-mongo.md) |
| ScyllaDB | [Test Result](./testing-result/testing-result-scylla.md) |

The [driver verification matrix](./testing.md#driver-verification-matrix) summarizes the scenarios and their status. All seven drivers passed the strict concurrency contract: all eight simultaneous valid creates returned HTTP 200, `success=true`, and matching persisted rows. MongoDB and ScyllaDB transaction rollback is not applicable to the GOREST ORM transaction API exercised here. The SQLite result page records the bundled SQLite engine and Go driver versions.

The default Go suite passed on October 23 with `go test ./... -count=1`, including app-level lifecycle checks. The driver test package skips its optional server integrations unless `GOREST_DRIVER_INTEGRATION` is set. These results describe the specific local server versions and configurations recorded; they do not establish compatibility with other versions, topologies, or deployment configurations. The run used source commit `1d86605f56cff6c5df751aabdb93bddc995b4870` with uncommitted working-tree changes.
