# Updating and validating the starter

Run the local quality checks before sharing or deploying changes:

```sh
make release-check
```

This runs the repository's formatting, static analysis, unit tests, lint, secret scan, and build checks. The required Go and tool versions are pinned in the repository.

To check Oracle behavior, install a compatible Oracle Instant Client, configure the dedicated non-production schema described in [Testing](./testing.md), set the exact schema confirmation, and run:

```sh
make test-integration
```

For combined test output and application coverage, use `make integration-reports`. That command runs all tagged integration tests; the Oracle suite writes to the confirmed schema, and optional Redis checks may connect to Redis if its configuration is enabled.

GitHub Actions runs unit tests and isolated Oracle integration-configuration checks. It does not connect to Oracle because the full integration suite requires an Oracle client, optional wallet, and a dedicated schema.

The starter currently uses GOREST v1.1.21 and the Go version in `go.mod`. Before changing framework or Go versions, review the target version's documentation and compatibility notes, then rerun the applicable checks and Oracle integration tests.
