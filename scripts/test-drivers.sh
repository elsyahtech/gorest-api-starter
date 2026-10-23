#!/usr/bin/env bash

set -Eeuo pipefail

repo_root=$(git rev-parse --show-toplevel 2>/dev/null) || {
	printf 'ERROR: Run this script from inside a Git repository.\n' >&2
	exit 1
}
cd "$repo_root"

compose_file="test/drivers/compose.yaml"
project="gorest-driver-it-$$"
current_service=''

cleanup() {
	local exit_code=$?
	trap - EXIT

	if [[ -n "$current_service" ]]; then
		docker compose --project-name "$project" --file "$compose_file" logs --no-color "$current_service" >&2 || true
	fi
	docker compose --project-name "$project" --file "$compose_file" down --remove-orphans || true
	exit "$exit_code"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
trap 'exit 129' HUP

if ! command -v docker >/dev/null 2>&1; then
	printf 'ERROR: Docker is required to run the full database integration suite.\n' >&2
	exit 1
fi

docker info >/dev/null

run_driver_test() {
	local driver=$1
	local service=${2:-}

	if [[ -n "$service" ]]; then
		current_service="$service"
		printf '\n=== Starting disposable %s service ===\n' "$driver"
		docker compose --project-name "$project" --file "$compose_file" up \
			--detach --wait --wait-timeout 360 "$service"
	fi

	printf '\n=== Running %s migration, seeder, CRUD, and transaction checks ===\n' "$driver"
	GOREST_DRIVER_INTEGRATION="$driver" go test -v -count=1 -timeout=20m ./test/drivers

	if [[ -n "$service" ]]; then
		docker compose --project-name "$project" --file "$compose_file" rm --stop --force "$service"
		current_service=''
	fi
}

run_oracle_test() {
	printf '\n=== Running Oracle checks with app/config and the configured local Instant Client ===\n'
	GOREST_DRIVER_INTEGRATION=oracle \
	GOREST_DRIVER_PORT="${GOREST_DRIVER_PORT:-1521}" \
		go test -v -count=1 -timeout=20m ./test/drivers
}

run_driver_test sqlite
run_driver_test postgres postgres
run_driver_test mysql mysql
run_oracle_test
run_driver_test sqlserver sqlserver
run_driver_test mongo mongo
run_driver_test scylla scylla

printf '\nAll seven database integration configurations passed.\n'
