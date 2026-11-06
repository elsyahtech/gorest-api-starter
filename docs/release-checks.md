# Updating and validating the starter

Run the local quality checks before sharing or deploying changes:

```sh
make release-check
```

This runs the repository's formatting, static analysis, unit tests, lint, secret scan, and build checks. The required Go and tool versions are pinned in the repository.

To check PostgreSQL behavior, configure the dedicated non-production database described in [Testing](./testing.md), then run:

```sh
make test-integration
```

For combined test output and application coverage, use `make integration-reports`. That command runs all tagged integration tests; optional Redis checks may connect to Redis if its configuration is enabled. Review the [Redis setup guidance](./configuration.md#redis-optional) before enabling it.

GitHub Actions runs the unit checks automatically. Its PostgreSQL integration job runs only when the required database and mTLS secrets are configured; without them, that job is skipped and no integration result is produced.

The starter currently uses GOREST v1.1.20 and the Go version in `go.mod`. Before changing framework or Go versions, review the target version's documentation and compatibility notes, then rerun the applicable checks and PostgreSQL integration tests.
