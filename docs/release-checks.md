# Release checks

The release-check workflow runs source formatting, `go vet`, unit tests, the baseline linter, a secret scan, and a build. The separate test workflow runs unit tests, coverage, and isolated ScyllaDB integration-configuration checks.

Run the local gate with the repository's pinned tools:

```sh
make release-check
```

To produce test evidence locally:

```sh
make coverage
make integration-reports
```

The integration report target runs the full tagged suite once and writes both report formats and application coverage. It requires the separately confirmed non-production ScyllaDB keyspace described in [Testing and coverage](./testing.md); optional Redis tests may run when enabled. GitHub Actions does not connect to external databases. Retain generated reports together with the source commit, Go version, ScyllaDB version, command, and sanitized environment description.

The Go toolchain version is `1.26.5`, as specified by `go.mod`. This starter currently targets GOREST v1.1.21. If your application needs a different GOREST version or connector capability, review the target version's [documentation and release history](https://github.com/elsyahtech/gorest/releases), assess configuration and API compatibility, and rerun the full release gate and ScyllaDB integration suite before adopting it.
