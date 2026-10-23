.PHONY: run lint lint-full test test-e2e test-mongo test-scylla test-drivers test-driver-sqlite coverage coverage-e2e format-check vet tool-versions secret-scan release-check check clean build ci upgrade fieldalignment rebuild push

run:
	@go run .

lint:
	@echo "Running pinned baseline linters..."
	golangci-lint run --enable-only=errcheck,govet,staticcheck,ineffassign,unused ./...

lint-full:
	golangci-lint run ./...

test:
	go test ./... -count=1

vet:
	go vet ./...

format-check:
	@files="$$(find . -name '*.go' -not -path './.git/*' -print0 | xargs -0 gofmt -l)"; \
	if [ -n "$$files" ]; then echo "Go files need gofmt:"; echo "$$files"; exit 1; fi

test-driver-sqlite:
	GOREST_DRIVER_INTEGRATION=sqlite go test ./test/drivers -run '^TestDatabaseDriverIntegration$$' -count=1 -v

secret-scan:
	gitleaks git . --redact --exit-code 1

tool-versions:
	@go version | grep -F 'go1.26.5' >/dev/null
	@golangci-lint --version | grep -F 'version 2.15.0' >/dev/null
	@go version -m "$$(command -v gitleaks)" | grep -E 'github.com/zricethezav/gitleaks/v8[[:space:]]+v8\.30\.1' >/dev/null

test-e2e:
	go test ./test/e2e -run E2E -count=1

test-mongo:
	GOREST_MONGO_E2E=1 go test ./test/e2e -run Mongo -count=1

test-scylla:
	GOREST_SCYLLA_E2E=1 go test ./test/e2e -run CQL -count=1

test-drivers:
	bash ./scripts/test-drivers.sh

coverage:
	go test -coverpkg=./app/... ./test/e2e -run E2E -count=1 -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

coverage-e2e:
	GOREST_MONGO_E2E=1 GOREST_SCYLLA_E2E=1 go test -coverpkg=./app/... ./test/e2e -run E2E -count=1 -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

check: vet test

build:
	mkdir -p bin
	go build -o bin/gorest-api-starter .

clean:
	rm -f bin/gorest-api-starter coverage.out

release-check: tool-versions format-check vet test lint test-driver-sqlite secret-scan build

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
