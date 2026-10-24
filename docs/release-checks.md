# Reproducible release checks

The `Release checks` workflow is the baseline release gate. It runs for pull requests, pushes to `main`, version tags, and manual dispatch. The Go version comes from `go.mod` (`1.26.5`); CI uses `GOTOOLCHAIN=local` so the installed toolchain must match that version.

## Run the gate locally

Install the pinned tools, then run the same Make target used in CI:

```sh
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0
go install github.com/zricethezav/gitleaks/v8@v8.30.1
GOTOOLCHAIN=local make release-check
```

The gate first checks the exact Go, golangci-lint, and Gitleaks versions. `make tool-versions` can run this check by itself and prints the required and detected versions, then reports whether each matches. The gate then runs Go formatting checks, `go vet`, uncached package tests, the pinned baseline linters (`errcheck`, `govet`, `staticcheck`, `ineffassign`, and `unused`), the SQLite E2E matrix target, Gitleaks against the working tree and full Git history, and a build. It stops with a non-zero status if a check fails. The Make targets are also available individually (`tool-versions`, `format-check`, `vet`, `test`, `lint`, `test-e2e-sqlite`, `secret-scan`, and `build`). The broader configured `make lint-full` remains available as a quality backlog check; it currently reports pre-existing findings in the starter's handlers, models, and formatting conventions.

## CI evidence

The workflow pins its GitHub Actions by commit SHA and installs the pinned Go analysis tools. For each run it uploads an artifact named with the source commit SHA. The artifact includes the check log and toolchain metadata; after a successful gate it also includes the Linux amd64 build and its SHA-256 checksum. GitHub retains the artifact for 90 days. Failed runs upload their log and metadata but do not include a release binary.

Use the workflow run URL, commit SHA, and downloaded artifact together as the release record. The binary checksum identifies the exact artifact that passed this gate. Rebuilds on another platform produce different artifacts and require their own checksums.

## Database driver evidence

The baseline CI gate runs the SQLite driver integration test because it does not need an external database service. It does not claim that PostgreSQL, MySQL, Oracle, SQL Server, MongoDB, or ScyllaDB passed in that CI run. Before a release claims one of those drivers as verified, run the selected integration test against a disposable service using the instructions in [Testing and E2E coverage](./testing.md#optional-database-driver-integration-tests). Record the server and tool versions, connection configuration without credentials, command, commit SHA, and complete test output on that driver's [evidence page](./testing-result.md).

Native MongoDB and ScyllaDB route scenarios run as part of `make test-e2e` and its driver-specific targets. Their coverage and container setup are documented in [Testing and E2E coverage](./testing.md#native-driver-routes).
