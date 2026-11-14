# Updating and validating the starter

Run the local quality checks before sharing or deploying changes:

```sh
make release-check
```

This runs the repository's formatting, static analysis, unit tests, lint, secret scan, and build checks. The required Go and tool versions are pinned in the repository.

To check MongoDB behavior, configure the dedicated non-production database described in [Testing](./testing.md), set the exact database confirmation, and run:

```sh
make test-integration
```

For combined test output and application coverage, use `GOREST_REDIS_DISABLE=true make integration-reports` to skip the regular Redis test. This does not disable the separate Redis mTLS test when its dedicated settings or runtime client certificate/key are configured. The MongoDB suite writes to the confirmed database; optional Redis tests may connect to Redis if enabled.

GitHub Actions runs unit tests and isolated MongoDB integration-configuration checks. It does not connect to MongoDB because integration tests require an explicitly managed server and dedicated test database.

The starter currently uses GOREST v1.1.21 and the Go version in `go.mod`. Before changing framework or Go versions, review the target version's documentation and compatibility notes, then rerun the applicable checks and MongoDB integration tests.
