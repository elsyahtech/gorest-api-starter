#!/usr/bin/env bash

set -Eeuo pipefail

repo_root=$(git rev-parse --show-toplevel 2>/dev/null) || {
	printf 'ERROR: Run this script from inside a Git repository.\n' >&2
	exit 1
}
cd "$repo_root"

compose_file="test/drivers/compose.yaml"
if [[ ! -f "$compose_file" ]]; then
	printf 'ERROR: Required Compose file is missing: %s\n' "$compose_file" >&2
	exit 1
fi

if ! command -v docker >/dev/null 2>&1 || ! docker compose version >/dev/null 2>&1; then
	printf 'ERROR: Docker Engine with Docker Compose v2 is required for Redis E2E tests.\n' >&2
	exit 1
fi
docker info >/dev/null

project="gorest-redis-e2e-$$"
service=redis
current_service=''

cleanup() {
	local exit_code=$?
	trap - EXIT
	if [[ -n "$current_service" ]]; then
		docker compose --project-name "$project" --file "$compose_file" logs --no-color "$current_service" >&2 || true
	fi
	docker compose --project-name "$project" --file "$compose_file" down --remove-orphans --volumes || true
	exit "$exit_code"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
trap 'exit 129' HUP

current_service=$service
printf '\n=== Starting disposable Redis service ===\n'
docker compose --project-name "$project" --file "$compose_file" up \
	--detach --wait --wait-timeout 60 "$service"
published_port=$(docker compose --project-name "$project" --file "$compose_file" port "$service" 6379)
published_port=${published_port##*:}

printf '\n=== Running Redis user CRUD and conflict checks ===\n'
GOREST_REDIS_E2E=1 \
GOREST_REDIS_HOST=127.0.0.1 \
GOREST_REDIS_PORT="$published_port" \
	go test -v -count=1 -timeout=5m ./test/e2e -run '^TestRedisUserCRUDE2E$'

docker compose --project-name "$project" --file "$compose_file" rm --stop --force "$service"
current_service=''
printf '\nAll Redis E2E scenarios passed.\n'