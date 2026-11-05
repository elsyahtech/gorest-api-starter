# Release checks

The release-check workflow runs source formatting, `go vet`, unit tests, the baseline linter, a secret scan, and a build. The separate test workflow runs unit tests and coverage, and runs MySQL E2E when the required external test database secrets are configured.

Run the local gate with the repository's pinned tools:

```sh
make release-check
```

To produce test evidence locally:

```sh
make coverage
make integration-reports
```

The integration report target runs the external MySQL suite once and writes both report formats and application coverage. In GitHub Actions, configure the required repository secrets and provision a reachable non-production schema before expecting E2E to run. When those secrets are absent, the E2E job is skipped; a skipped job is not a passing integration test. See [Testing and coverage](./testing.md#github-actions-ci) for CI setup and [MySQL configuration](./database.md) for database requirements. Retain generated reports together with the source commit, Go version, MySQL version, command, and sanitized environment description.

The Go toolchain version is `1.26.5`, as specified by `go.mod`. This starter currently targets GOREST v1.1.20. If your application needs a different GOREST version or connector capability, review the target version's [documentation and release history](https://github.com/elsyahtech/gorest/releases), assess configuration and API compatibility, and rerun the full release gate and MySQL integration suite before adopting it.
