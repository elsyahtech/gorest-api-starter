# Testing and coverage

This guide describes the test strategy for the active MySQL user API. Unit tests are isolated and do not connect to a database. E2E tests build and start the application, send HTTP requests, and use only a separately configured external MySQL test database. No Docker or local disposable database is used.

## Test layout

```text
app/
  config/*_test.go          configuration and validation tests
  handler/*_test.go         HTTP handler unit tests with a stub use case
  library/*_test.go         user business-logic unit tests with a stub repository
  model/*_test.go           model input guards
  router/*_test.go          route registration and authentication wiring
  view/*_test.go            entity-to-response mapping tests
integration-test/
  config.go                 test-only environment and database-safety checks
  users_http_test.go        HTTP E2E flow against the running application
```

The current code has no separate service package or application-owned middleware package requiring independent unit tests. GOREST's authentication filter is verified through the HTTP integration flow.

## Unit tests

Run without MySQL:

```sh
make unit-test
# equivalent:
go test ./... -count=1
```

Examples of isolated unit checks include password hashing before repository persistence, rejected invalid model inputs, handler request/response behavior with stubs, and public response mapping that omits password fields. The tests do not claim to validate MySQL connectivity or SQL behavior.

## Configure the external MySQL integration test

1. Prepare a dedicated MySQL test database on the host/IP and port to be tested. The schema name must end with `_test` or `_testing`; for example, `gorest_starter_test`.
2. Apply the MySQL migration to that test schema and ensure the `users` table exists. The test does not run migrations, seed data, or drop tables.
3. Use the existing `.env` copied from `.env.example` and fill in the `GOREST_TEST_MYSQL_*` fields there. Alternatively, export those variables in the shell. Do not create a second integration environment file. Missing test connection values fall back to `GOREST_DB_*` only when `GOREST_APP_ENV` is `development`, `test`, or `testing`, the driver is MySQL, and the database is enabled. Prefer explicit test connection values to keep the test target unambiguous.
4. Set `GOREST_TEST_MYSQL_CONFIRM_DATABASE` to exactly the configured test database name. This is a required explicit opt-in before test data can be written.
5. Ensure the test account can insert, read, and delete only the intended test data.

The test uses `GOREST_TEST_MYSQL_HOST`, `PORT`, `DATABASE`, `USER`, `PASSWORD`, and `TLS`, documented in `.env.example`. TLS uses Go/system trusted roots and does not require a CA path. For a private CA, configure `SSL_CERT_FILE` to a PEM bundle containing the system roots and the private CA. The test runner maps validated connection settings to the app's runtime variables for its child process. The harness closes idle preflight/application database connections so E2E can run under accounts with constrained connection limits. Do not set these variables to production credentials or a production database. The test inserts a UUID-scoped user and cleanup deletes only that exact ID/email record. It does not truncate or drop tables.

Database names normally must end in `_test` or `_testing`. If a dedicated non-production test database uses another name, an additional explicit opt-in is required: set `GOREST_TEST_MYSQL_ALLOW_NON_TEST_DATABASE` to exactly the same database name as the confirmation variable. Never use this override for production.

Run the integration suite separately:

```sh
make test-integration
# equivalent:
go test -tags=integration -v -count=1 -timeout=5m ./integration-test
```

The Go build tag keeps external-database tests out of `go test ./...`. If required test settings or the safety confirmation are missing, the integration run fails before starting the application or modifying test data.

For a dedicated non-production test database whose name does not end in `_test` or `_testing`, explicitly confirm the exact name in both variables:

```sh
GOREST_TEST_MYSQL_CONFIRM_DATABASE='your-test-database' \
GOREST_TEST_MYSQL_ALLOW_NON_TEST_DATABASE='your-test-database' \
make test-integration
```

If the host resolves to an unreachable IPv6 address in your environment, run with `GODEBUG=netdns=cgo` to use the system resolver.

## GitHub Actions CI

`.github/workflows/test.yml` runs unit tests and generates coverage on pushes, pull requests, and manual dispatches. It uploads `coverage/unit.coverprofile` to Codecov and retains the unit coverage reports as a workflow artifact. The Codecov upload is informational and does not fail CI if upload fails.

To run the HTTP integration job in GitHub Actions, configure these repository Actions secrets: `GOREST_TEST_MYSQL_HOST`, `GOREST_TEST_MYSQL_DATABASE`, `GOREST_TEST_MYSQL_USER`, `GOREST_TEST_MYSQL_PASSWORD`, and `GOREST_TEST_MYSQL_CONFIRM_DATABASE`. The confirmation must exactly match the test database name. Database names must end in `_test` or `_testing`, unless `GOREST_TEST_MYSQL_ALLOW_NON_TEST_DATABASE` is also configured with the exact same name. Set the optional `GOREST_TEST_MYSQL_TLS` secret to `true` when the server requires TLS; a CA path is not required when the server certificate chains to a runner system root. For a private CA, store the PEM CA certificate contents in the optional `GOREST_TEST_MYSQL_CA_PEM` secret. The workflow adds it to the runner's system trust bundle before running the tests.

Before enabling the job, provision a dedicated non-production database and its `users` table, verify that GitHub-hosted runners can reach the database, and grant the test account only the permissions needed by the test. The test does not create schema or tables. Store credentials and the PEM CA contents as GitHub Actions secrets rather than workflow literals; see GitHub's [using secrets in GitHub Actions](https://docs.github.com/en/actions/security-for-github-actions/security-guides/using-secrets-in-github-actions) guide. GOREST v1.1.19 supports TLS with system roots and hostname verification; see [MySQL TLS configuration](./database.md#configure-tls-private-cas-and-client-certificates) for private CAs and the mTLS limitation.

If the required secrets are unavailable, GitHub skips the E2E job; it does not substitute a local or production database. A skipped job means E2E was not run and is not evidence that the HTTP flow passed. This commonly applies to pull requests from forks, for which GitHub does not expose repository secrets to the workflow. See GitHub's [events that trigger workflows](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#pull_request) documentation for the `pull_request` event's secret-access behavior. Configure repository Actions secrets for a dedicated non-production database according to your project's access policy; do not weaken the workflow's database safety checks to make secrets available.

The Codecov token is read from the optional `CODECOV_TOKEN` repository secret. The uploaded file is the unit profile generated by `make coverage`; E2E runs the application as a separate process and does not produce an application integration-coverage profile.

## HTTP scenarios

`TestMySQLUserEndpointsOverHTTP` checks:

- The running application becomes ready.
- Malformed JSON for `POST /users/create` returns an error response.
- Missing required user fields and duplicate user creation are rejected.
- A valid create request returns success, inserts a MySQL row, and stores a hash instead of the submitted password.
- Incorrect login credentials return HTTP `401`; a valid login returns HTTP `200` and an `access` cookie.
- `GET /users/` without authentication returns HTTP `401`.
- `GET /users/` with the access cookie returns the created user without a password field.
- Test data is removed by a targeted ID/email cleanup.

The test uses a random user ID/email on each run so repeated and concurrent runs do not collide. The suite does not mutate schema and requires the `users` table to be provisioned before execution.

## Coverage and reports

Generate the unit coverage profile, function summary, and HTML report:

```sh
make coverage
```

Generated files:

- `coverage/unit.coverprofile` — Go coverage profile for machine processing.
- `coverage/unit-coverage.txt` — human-readable per-function and total summary from `go tool cover -func`.
- `coverage/unit-coverage.html` — source-annotated HTML report from `go tool cover -html`.

The coverage command instruments `./app/...` and runs the package unit tests. With the current test harness, integration E2E executes the application as a separately built process, so the unit profile does **not** count application statements exercised in that child process. The E2E report is execution evidence, not an integration coverage percentage.

If your project requires coverage that includes application code exercised by end-to-end tests, use Go's supported integration-coverage workflow: build the application with coverage instrumentation, set a dedicated `GOCOVERDIR` for the instrumented process, and merge/read the resulting data with `go tool covdata`. See the Go team's [Code coverage for Go integration tests](https://go.dev/doc/build-cover) guide. Treat this as a separate test-pipeline change: preserve the existing unit profile, ensure child processes terminate cleanly so coverage data is flushed, and validate the combined result against a known scenario.

Generate E2E reports after configuring the external test database:

```sh
make integration-report # human-readable Go test output: coverage/integration-test.txt
make integration-json   # machine-readable Go test events: coverage/integration-test.json
```

Each target reruns the integration suite. Reports under `coverage/` are generated artifacts and are ignored by Git; regenerate them for each source revision and test environment whose results you want to retain. For reproducible test evidence, keep the source commit, Go version, MySQL server version, test command, sanitized configuration (no credentials), and report filenames with the reports. A coverage percentage alone does not demonstrate that the application has been tested sufficiently.

## Reading the evidence

- A passing unit suite means the isolated package tests passed on the current source.
- A passing E2E suite means the listed HTTP workflow passed against the MySQL server configured for that run.
- The coverage total is based only on the instrumented unit-test run described above.
- A test run is evidence for TLS/mTLS only when the configured MySQL connector verifies the required server certificate and, for mTLS, presents an accepted client certificate. Follow [MySQL TLS configuration](./database.md#configure-tls-private-cas-and-client-certificates) and verify both accepted and rejected connection cases.
- Keep each report with its source revision and environment details so another team member can understand and reproduce the test run.
