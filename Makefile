SHELL := /bin/bash

.PHONY: run test unit-test test-integration integration-report integration-json coverage vet format-check lint lint-full tool-versions secret-scan release-check check clean build ci upgrade fieldalignment rebuild push

run:
	go run .

test: unit-test

unit-test:
	go test -v ./... -count=1

test-integration:
	go test -tags=integration -v -count=1 -timeout=5m ./integration-test

integration-report:
	mkdir -p coverage
	set -o pipefail; go test -tags=integration -v -count=1 -timeout=5m ./integration-test 2>&1 | tee coverage/integration-test.txt

integration-json:
	mkdir -p coverage
	go test -tags=integration -json -count=1 -timeout=5m ./integration-test > coverage/integration-test.json

coverage:
	mkdir -p coverage
	go test -coverpkg=./app/... -coverprofile=coverage/unit.coverprofile ./app/...
	go tool cover -func=coverage/unit.coverprofile | tee coverage/unit-coverage.txt
	go tool cover -html=coverage/unit.coverprofile -o coverage/unit-coverage.html

vet:
	go vet ./...

format-check:
	@files="$$(find . -name '*.go' -not -path './.git/*' -print0 | xargs -0 gofmt -l)"; \
	if [ -n "$$files" ]; then echo "Go files need gofmt:"; echo "$$files"; exit 1; fi

lint:
	golangci-lint run --enable-only=errcheck,govet,staticcheck,ineffassign,unused ./...

lint-full:
	golangci-lint run ./...

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
	case "$$lint_actual" in *"version 2.15.0"*) echo "PASS: $$lint_actual" ;; *) echo "FAIL: expected golangci-lint 2.15.0; found: $$lint_actual"; exit 1 ;; esac; \
	if ! command -v gitleaks >/dev/null 2>&1; then echo "FAIL: Gitleaks is not installed or not on PATH (required: 8.30.1)."; exit 1; fi; \
	gitleaks_bin=$$(command -v gitleaks); \
	gitleaks_actual=$$(go version -m "$$gitleaks_bin" | awk '$$1 == "mod" && $$2 == "github.com/zricethezav/gitleaks/v8" { print $$3 }'); \
	if [ "$$gitleaks_actual" != "v8.30.1" ]; then echo "FAIL: expected Gitleaks v8.30.1; found: $${gitleaks_actual:-unknown build version}"; exit 1; fi; \
	echo "PASS: Gitleaks $$gitleaks_actual"; \
	echo "All pinned tool versions match."

secret-scan:
	gitleaks dir . --redact --exit-code 1
	gitleaks git . --redact --exit-code 1

release-check: tool-versions format-check vet unit-test lint secret-scan build

check: vet unit-test

build:
	mkdir -p bin
	go build -o bin/gorest-api-starter .

ci: release-check

upgrade:
	go get -u ./...
	go mod tidy

fieldalignment:
	fieldalignment -fix ./...

rebuild:
	mkdir -p bin
	go build -a -o bin/gorest-api-starter .

push:
	bash ./scripts/push.sh

clean:
	rm -f bin/gorest-api-starter coverage/unit.coverprofile coverage/unit-coverage.txt coverage/unit-coverage.html coverage/integration-test.txt coverage/integration-test.json
