.PHONY: run lint lint-full test test-e2e test-e2e-sqlite test-e2e-postgres test-e2e-mysql test-e2e-oracle test-e2e-sqlserver test-e2e-mongo test-e2e-scylla test-e2e-redis coverage coverage-e2e format-check vet tool-versions secret-scan release-check check clean build ci upgrade fieldalignment rebuild push

run:
	@go run .

lint:
	@echo "Running pinned baseline linters..."
	golangci-lint run --enable-only=errcheck,govet,staticcheck,ineffassign,unused ./...

lint-full:
	golangci-lint run ./...

test:
	@echo "Running all Go tests (uncached; verbose output)..."
	go test -v ./... -count=1
	@echo "All Go tests passed."

vet:
	@echo "Running go vet across all packages..."
	go vet ./...
	@echo "go vet passed for all packages."

format-check:
	@echo "Checking Go formatting..."
	@files="$$(find . -name '*.go' -not -path './.git/*' -print0 | xargs -0 gofmt -l)"; \
	if [ -n "$$files" ]; then echo "Go files need gofmt:"; echo "$$files"; exit 1; fi; \
	echo "All Go files are gofmt-formatted."

secret-scan:
	gitleaks dir . --redact --exit-code 1
	gitleaks git . --redact --exit-code 1

tool-versions:
	@set -eu; \
	echo "Checking pinned tool versions required by the release gate..."; \
	if ! command -v go >/dev/null 2>&1; then echo "FAIL: Go is not installed or not on PATH (required: go1.26.5)."; exit 1; fi; \
	go_actual=$$(go version | awk '{print $$3}'); \
	echo "Checking Go (required: go1.26.5; found: $$go_actual)..."; \
	if [ "$$go_actual" != "go1.26.5" ]; then echo "FAIL: install Go 1.26.5."; exit 1; fi; \
	echo "PASS: Go $$go_actual"; \
	if ! command -v golangci-lint >/dev/null 2>&1; then echo "FAIL: golangci-lint is not installed or not on PATH (required: 2.15.0)."; exit 1; fi; \
	lint_actual=$$(golangci-lint --version); \
	echo "Checking golangci-lint (required: 2.15.0)..."; \
	case "$$lint_actual" in *"version 2.15.0"*) echo "PASS: $$lint_actual" ;; *) echo "FAIL: expected golangci-lint 2.15.0; found: $$lint_actual"; exit 1 ;; esac; \
	if ! command -v gitleaks >/dev/null 2>&1; then echo "FAIL: Gitleaks is not installed or not on PATH (required: 8.30.1)."; exit 1; fi; \
	gitleaks_bin=$$(command -v gitleaks); \
	gitleaks_actual=$$(go version -m "$$gitleaks_bin" | awk '$$1 == "mod" && $$2 == "github.com/zricethezav/gitleaks/v8" { print $$3 }'); \
	echo "Checking Gitleaks (required: v8.30.1; binary: $$gitleaks_bin)..."; \
	if [ "$$gitleaks_actual" != "v8.30.1" ]; then echo "FAIL: expected Gitleaks v8.30.1; found: $${gitleaks_actual:-unknown build version}"; exit 1; fi; \
	echo "PASS: Gitleaks $$gitleaks_actual"; \
	echo "All pinned tool versions match."

test-e2e:
	@echo "Running the complete isolated E2E matrix for all seven database drivers..."
	bash ./scripts/test-e2e.sh

test-e2e-sqlite:
	bash ./scripts/test-e2e.sh sqlite

test-e2e-postgres:
	bash ./scripts/test-e2e.sh postgres

test-e2e-mysql:
	bash ./scripts/test-e2e.sh mysql

test-e2e-oracle:
	bash ./scripts/test-e2e.sh oracle

test-e2e-sqlserver:
	bash ./scripts/test-e2e.sh sqlserver

test-e2e-mongo:
	bash ./scripts/test-e2e.sh mongo

test-e2e-scylla:
	bash ./scripts/test-e2e.sh scylla

test-e2e-redis:
	bash ./scripts/test-e2e-redis.sh

coverage:
	go test -coverpkg=./app/... ./test/e2e -run E2E -count=1 -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

coverage-e2e:
	GOREST_MONGO_E2E=1 GOREST_SCYLLA_E2E=1 go test -coverpkg=./app/... ./test/e2e -run E2E -count=1 -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

check:
	@echo "Running static analysis and the full Go test suite..."
	@$(MAKE) vet
	@$(MAKE) test
	@echo "All checks passed."

build:
	mkdir -p bin
	go build -o bin/gorest-api-starter .

clean:
	rm -f bin/gorest-api-starter coverage.out

release-check: tool-versions format-check vet test lint test-e2e-sqlite secret-scan build

ci: release-check

upgrade:
	@echo "Upgrading dependencies..."
	go get -u ./...
	go mod tidy

fieldalignment:
	fieldalignment -fix ./...

rebuild:
	mkdir -p bin
	go build -a -o bin/gorest-api-starter .

push:
	bash ./scripts/push.sh
