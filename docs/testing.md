# Testing and coverage

Unit tests do not require a ScyllaDB server. The tagged database integration tests connect to a separately configured external ScyllaDB test keyspace and exercise migrations, seeded data, library calls, and HTTP behavior. No Docker or disposable local database is used.

## Test layout

```text
app/
  config/*_test.go          application configuration and validation
  library/*_test.go         isolated input validation
  model/*_test.go           model input guards
  router/*_test.go          route registration and authentication wiring
  view/*_test.go            entity-to-response mapping
integration-test/
  config.go                 ScyllaDB test configuration and safety checks
  config_test.go            configuration safety unit tests
  scylla_connection_test.go external ScyllaDB connection test
  scylla_mtls_test.go       optional client-certificate connection test
  user_library_test.go      user library behavior against ScyllaDB
  users_http_test.go        HTTP E2E flow against ScyllaDB
  hardening_http_test.go    HTTP hardening-report integration test
  redis*_test.go            optional external Redis checks
scripts/
  integration-report.go     single-run integration reports and app coverage
```

The hardening-report test runs independently without a database:

```sh
make test-hardening-integration
```

## Unit tests

```sh
make unit-test
# equivalent:
go test ./... -count=1
```

Unit tests cover configuration validation, input guards, route/authentication wiring, and public response mapping. ScyllaDB-dependent persistence and HTTP behavior are covered by the tagged integration tests.

## Configure the ScyllaDB integration tests

1. Provision a dedicated non-production ScyllaDB keyspace. It must already exist before the test connects; a name ending in `_test` or `_testing` is recommended.
2. Set `GOREST_TEST_SCYLLA_HOST`, `PORT`, `KEYSPACE`, and, when needed, `USER` and `PASSWORD` in `.env` or the process environment. The optional TLS settings are `GOREST_TEST_SCYLLA_TLS`, `CA_FILE`, `CERT_FILE`, and `KEY_FILE`; use `enable` for verified TLS and configure the client certificate and key together.
3. Set `GOREST_TEST_SCYLLA_CONFIRM_KEYSPACE` to exactly match the configured test keyspace. For a keyspace without the `_test` or `_testing` suffix, also set `GOREST_TEST_SCYLLA_ALLOW_NON_TEST_KEYSPACE` to that exact name.
4. Ensure the test account can connect, create the users table and email index, use GOREST migration lock tables/lightweight transactions, and perform the test's insert/read/delete operations. Verify the needed permissions with your ScyllaDB administrator.

If the dedicated test host/keyspace/credentials are not configured, the test harness can use `GOREST_DB_*` only when `GOREST_APP_ENV` is `development`, `test`, or `testing`, `GOREST_DB_DRIVER=scylla`, and `GOREST_DB_DISABLE=false`. Exact keyspace confirmation is still required. Test-specific TLS paths override application paths. The test runner maps the validated settings into the migration and HTTP child process; do not set them to production credentials or a production keyspace.

The library and HTTP tests run the up migration and seeder. They do not drop the keyspace or tables and leave migration/seeder metadata and ten sample rows in the test keyspace. Each test uses a unique user ID and deletes only its own records. Use a disposable keyspace or review the retained artifacts after testing.

Run the selected ScyllaDB integration tests:

```sh
make test-integration
# run the full tagged suite, including optional Redis checks:
make integration-reports
```

The Go `integration` build tag excludes database-dependent tests from `go test ./...`. `make test-integration` fails if the ScyllaDB configuration or exact keyspace confirmation is missing. The mTLS check is skipped when client-certificate settings are absent and fails on partial configuration. The full tagged suite may connect to Redis if Redis is enabled in the environment; leave Redis disabled to avoid those optional checks.

## GitHub Actions CI

The test workflow runs unit tests and isolated ScyllaDB integration-configuration checks. It does not connect to ScyllaDB or Redis. External database tests are run locally or in a separately managed environment with explicit test-keyspace confirmation. A passing CI workflow does not claim that external persistence, TLS, or HTTP E2E tests passed.

The release-check workflow runs formatting, `go vet`, unit tests, linting, secret scanning, and a build. See [Release checks](./release-checks.md).

## HTTP scenarios

`TestScyllaUserEndpointsOverHTTP` checks application startup against ScyllaDB and verifies:

- Malformed JSON and missing required user fields are rejected.
- Duplicate user creation is rejected.
- A valid create request stores a password hash rather than the submitted password.
- Incorrect login credentials return HTTP `401`; valid credentials return HTTP `200` and an access cookie.
- `GET /users/` without authentication returns HTTP `401`.
- `GET /users/` with the access cookie returns the created user without a password field.
- Cleanup deletes only the unique ID/email record created by the test.

## Optional TLS and Redis checks

`TestScyllaMutualTLSRequiresClientCertificate` is an opt-in connection test using the ScyllaDB test TLS settings; configure an endpoint that requires client certificates. Keep private keys out of Git and do not point the test at production.

The Redis integration tests are independent of the current user routes. They run only when the corresponding Redis service settings are enabled/configured. They use unique test keys and do not flush Redis; the Redis mTLS test requires an endpoint configured to require client certificates. Keep optional Redis settings disabled unless they point to dedicated non-production test infrastructure.

## Coverage and reports

Generate the unit profile, function summary, and HTML report:

```sh
make coverage
```

To run the full tagged integration suite once and produce both readable/JSON test reports and application coverage:

```sh
make integration-reports
```

Outputs are generated under `coverage/` and ignored by Git:

- `unit.coverprofile`, `unit-coverage.txt`, and `unit-coverage.html`
- `integration-test.txt` and `integration-test.json`
- `integration.coverprofile`, `integration-coverage.txt`, and `integration-coverage.html`

Keep generated evidence with the source commit, Go version, ScyllaDB server version, test command, and sanitized environment description. A coverage percentage alone does not demonstrate that the application has been tested sufficiently. See the Go team's [Code coverage for Go integration tests](https://go.dev/doc/build-cover) guide.

## Reading the evidence

- Passing unit tests means isolated package tests passed on the current source.
- Passing ScyllaDB integration tests means the covered flows passed against the configured keyspace for that run.
- A test run is evidence for TLS/mTLS only when the endpoint uses the expected server certificate and client-certificate policy.
- Unit and integration profiles cover different execution paths; retain both with their source revision and environment details.
