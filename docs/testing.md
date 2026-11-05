# Testing and coverage

This guide describes the test strategy for the active MySQL user API and the GOREST services it uses. Unit tests are isolated and do not connect to external services. HTTP E2E tests build and start the application, send requests, and use a separately configured external MySQL test database. Optional direct integration tests verify MySQL mTLS and Redis against explicitly configured external test services. No Docker or local disposable database is used.

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
  mysql_mtls_test.go        optional MySQL client-certificate verification
  redis_integration_test.go optional real Redis service test
  redis_mtls_test.go        optional Redis client-certificate verification
scripts/
  integration-report.go     one-run test reporting and application coverage
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

The HTTP test uses `GOREST_TEST_MYSQL_HOST`, `PORT`, `DATABASE`, `USER`, `PASSWORD`, and `TLS`, documented in `.env.example`. It also accepts optional `GOREST_TEST_MYSQL_CA_FILE`, `CERT_FILE`, and `KEY_FILE` paths for a private server CA and mTLS. The client certificate and key must be configured together. If those test-specific paths are blank, the test inherits the matching `GOREST_DB_*` CA/certificate settings when using the application database as its test endpoint. The test runner maps the validated connection settings to the preflight connection and app child process. The harness closes idle preflight/application database connections so E2E can run under accounts with constrained connection limits. Do not set these variables to production credentials or a production database. The test inserts a UUID-scoped user and cleanup deletes only that exact ID/email record. It does not truncate or drop tables.

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

To also run the mTLS connection test in CI, configure all of the optional `GOREST_TEST_MYSQL_MTLS_HOST`, `PORT`, `DATABASE`, `USER`, `PASSWORD`, `CA_PEM`, `CLIENT_CERT_PEM`, and `CLIENT_KEY_PEM` secrets. The workflow writes certificate/key contents to restricted-permission temporary files and maps their paths for the test. The dedicated MySQL account must require client certificates (`REQUIRE X509`); mTLS test secrets are independent of the HTTP E2E database settings.

To run the Redis integration test in CI, configure the optional `GOREST_TEST_REDIS_HOST`, `PORT`, `DATABASE`, `USER`, `PASSWORD`, and `TLS` secrets. `GOREST_TEST_REDIS_CA_PEM` can provide a private CA; `GOREST_TEST_REDIS_CLIENT_CERT_PEM` and `GOREST_TEST_REDIS_CLIENT_KEY_PEM` are optional and must be configured together. Use a dedicated Redis test instance/database and a limited ACL account.

To run the Redis mTLS check in CI, configure the optional `GOREST_TEST_REDIS_MTLS_HOST`, `DATABASE`, `CA_PEM`, `CLIENT_CERT_PEM`, and `CLIENT_KEY_PEM` secrets; `PORT` defaults to `6379`, and `USER` / `PASSWORD` are needed only when Redis authentication is enabled. This must be an endpoint configured to require client certificates. The workflow writes the CA and client key pair to restricted-permission temporary files. See [Redis TLS and mTLS configuration](./configuration.md#use-a-private-ca-or-mutual-tls) for client-side setup and certificate roles.

Before enabling the job, provision a dedicated non-production database and its `users` table, verify that GitHub-hosted runners can reach the database, and grant the test account only the permissions needed by the test. The test does not create schema or tables. Store credentials and the PEM CA contents as GitHub Actions secrets rather than workflow literals; see GitHub's [using secrets in GitHub Actions](https://docs.github.com/en/actions/security-for-github-actions/security-guides/using-secrets-in-github-actions) guide. GOREST v1.1.20 supports TLS with system roots or a configured CA and can present a configured client certificate for mTLS; see [MySQL TLS configuration](./database.md#configure-tls-private-cas-and-client-certificates).

If the required secrets are unavailable, GitHub skips the E2E job; it does not substitute a local or production database. A skipped job means E2E was not run and is not evidence that the HTTP flow passed. This commonly applies to pull requests from forks, for which GitHub does not expose repository secrets to the workflow. See GitHub's [events that trigger workflows](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#pull_request) documentation for the `pull_request` event's secret-access behavior. Configure repository Actions secrets for a dedicated non-production database according to your project's access policy; do not weaken the workflow's database safety checks to make secrets available.

The Codecov token is read from the optional `CODECOV_TOKEN` repository secret. CI uploads the unit profile from `make coverage` and, when the external database is configured and the integration job succeeds, the application profile from `make integration-reports`. If mTLS or Redis settings are configured, those optional tests run as part of the same suite; without their settings, they are explicitly skipped.

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

## MySQL mutual-TLS integration test

`TestMySQLMutualTLSRequiresClientCertificate` is an opt-in, read-only connection test using GOREST v1.1.20's native MySQL TLS paths. It is separate from the HTTP E2E test. To run it, configure all of the following in `.env` or the environment:

- `GOREST_TEST_MYSQL_MTLS_HOST`, optional `GOREST_TEST_MYSQL_MTLS_PORT` (defaults to `3306`), and `GOREST_TEST_MYSQL_MTLS_DATABASE`.
- `GOREST_TEST_MYSQL_MTLS_USER` and `GOREST_TEST_MYSQL_MTLS_PASSWORD` for a dedicated MySQL account that has `REQUIRE X509`.
- `GOREST_TEST_MYSQL_MTLS_TLS=true`.
- `GOREST_TEST_MYSQL_MTLS_CA_FILE`, `GOREST_TEST_MYSQL_MTLS_CLIENT_CERT_FILE`, and `GOREST_TEST_MYSQL_MTLS_CLIENT_KEY_FILE`, pointing to readable PEM files.

Alternatively, in a development/test environment, the test can reuse `GOREST_DB_HOST`, `PORT`, `NAME`, `USER`, `PASSWORD`, `TLS`, `CA_PATH`, `CERT_PATH`, and `KEY_PATH`. Explicit `GOREST_TEST_MYSQL_MTLS_*` values take precedence. Run the test with `go test -tags=integration -run '^TestMySQLMutualTLSRequiresClientCertificate$' -count=1 ./integration-test`. The server must reject the account without a client certificate (error 1045) and accept it with the configured certificate; the test also checks that the accepted session has a non-empty `Ssl_cipher`. If no mTLS configuration is present, the test is skipped; partial configuration is an error. Keep private keys out of Git and never point the test at production.

## Redis integration test

`TestGORESTRedisServiceAgainstConfiguredInstance` uses the Redis connection settings from the process environment or `.env`. Set `GOREST_REDIS_DISABLE=false` only when those settings point to a dedicated non-production Redis instance. The test refuses `GOREST_APP_ENV=production`, checks connectivity, and reads, writes, and deletes only a UUID-scoped key. It does not flush Redis or modify existing keys. Use an account with permission for the selected test database and key operations.

For TLS server authentication, set `GOREST_REDIS_TLS=true`. A Redis certificate signed by a system-trusted CA does not require `GOREST_REDIS_CA_PATH`; use that setting for a private CA. The optional `TestGORESTRedisMutualTLSRequiresClientCertificate` checks a Redis endpoint configured to require client certificates: the endpoint must reject the connection without a client identity and accept it with the configured certificate/key. Configure a CA file and both client certificate files using either `GOREST_TEST_REDIS_MTLS_*` for a separate endpoint, or `GOREST_REDIS_CA_PATH`, `GOREST_REDIS_CERT_PATH`, and `GOREST_REDIS_KEY_PATH` with `GOREST_REDIS_TLS=true` for the runtime endpoint. The test skips when no client certificate is configured and rejects a partial certificate/key pair. Do not use production credentials or services.

## Coverage and reports

Generate the unit coverage profile, function summary, and HTML report:

```sh
make coverage
```

Generated files:

- `coverage/unit.coverprofile` — Go coverage profile for machine processing.
- `coverage/unit-coverage.txt` — human-readable per-function and total summary from `go tool cover -func`.
- `coverage/unit-coverage.html` — source-annotated HTML report from `go tool cover -html`.

The unit coverage command instruments `./app/...` and runs the package unit tests. Generate E2E reports and application integration coverage together after configuring the external test database:

```sh
make integration-reports
```

This runs the integration suite once. It instruments the application child process, waits for its graceful shutdown to flush coverage, and generates:

- `coverage/integration-test.txt` — readable Go test output.
- `coverage/integration-test.json` — machine-readable Go test events.
- `coverage/integration.coverprofile` — coverage profile for the application entry point and `./app/...` exercised by E2E.
- `coverage/integration-coverage.txt` and `coverage/integration-coverage.html` — function summary and source-annotated report.

`make integration-report` and `make integration-json` remain aliases for the same combined target. Reports under `coverage/` are generated artifacts and are ignored by Git; regenerate them for each source revision and test environment whose results you want to retain. For reproducible test evidence, keep the source commit, Go version, MySQL server version, test command, sanitized configuration (no credentials), and report filenames with the reports. A coverage percentage alone does not demonstrate that the application has been tested sufficiently. See the Go team's [Code coverage for Go integration tests](https://go.dev/doc/build-cover) guide.

## Reading the evidence

- A passing unit suite means the isolated package tests passed on the current source.
- A passing E2E suite means the listed HTTP workflow passed against the MySQL server configured for that run.
- The unit profile and integration profile are generated by separate test runs and cover different execution paths.
- A test run is evidence for TLS/mTLS only when the configured MySQL connector verifies the required server certificate and, for mTLS, presents an accepted client certificate. Follow [MySQL TLS configuration](./database.md#configure-tls-private-cas-and-client-certificates) and verify both accepted and rejected connection cases.
- Keep each report with its source revision and environment details so another team member can understand and reproduce the test run.
