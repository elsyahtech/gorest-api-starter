.PHONY: run lint test test-e2e test-mongo test-scylla coverage coverage-e2e check clean build ci upgrade fieldalignment rebuild push

run:
	@go run .

lint:
	@echo "Running only golangci-lint..."
	golangci-lint run ./...

test:
	go test ./...

test-e2e:
	go test ./test/e2e -run E2E -count=1

test-mongo:
	GOREST_MONGO_E2E=1 go test ./test/e2e -run Mongo -count=1

test-scylla:
	GOREST_SCYLLA_E2E=1 go test ./test/e2e -run CQL -count=1

coverage:
	go test -coverpkg=./app/... ./test/e2e -run E2E -count=1 -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

coverage-e2e:
	GOREST_MONGO_E2E=1 GOREST_SCYLLA_E2E=1 go test -coverpkg=./app/... ./test/e2e -run E2E -count=1 -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

check:
	go vet ./...
	go test ./...

build:
	mkdir -p bin
	go build -o bin/gorest-api-starter .

clean:
	rm -f bin/gorest-api-starter coverage.out

ci: check lint

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