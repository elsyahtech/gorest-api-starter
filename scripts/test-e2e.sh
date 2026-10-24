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

if [[ $# -gt 1 ]]; then
	printf 'Usage: %s [sqlite|postgres|mysql|oracle|sqlserver|mongo|scylla]\n' "$0" >&2
	exit 2
fi
selected_driver=${1:-all}
case "$selected_driver" in
	all|sqlite|postgres|mysql|oracle|sqlserver|mongo|scylla) ;;
	*) printf 'ERROR: unsupported driver %q. Choose sqlite, postgres, mysql, oracle, sqlserver, mongo, or scylla.\n' "$selected_driver" >&2; exit 2 ;;
esac

case "$(uname -s)" in
	Darwin|Linux) ;;
	*)
		printf 'ERROR: Automated driver testing currently supports macOS and Linux only.\n' >&2
		exit 1
		;;
esac

docker_enabled=0
if [[ "$selected_driver" != sqlite ]]; then
	if ! command -v docker >/dev/null 2>&1; then
		printf 'ERROR: Docker is required to run the database integration suite.\n' >&2
		exit 1
	fi
	if ! docker compose version >/dev/null 2>&1; then
		printf 'ERROR: Docker Compose v2 is required.\n' >&2
		exit 1
	fi
	docker info >/dev/null
	docker_enabled=1
fi

project="gorest-driver-it-$$"
temporary_root=$(mktemp -d "${TMPDIR:-/tmp}/gorest-driver-it.XXXXXX")
current_service=''
oracle_mount=''

cleanup() {
	local exit_code=$?
	trap - EXIT

	if [[ -n "$current_service" ]]; then
		docker compose --project-name "$project" --file "$compose_file" logs --no-color "$current_service" >&2 || true
	fi
	if [[ "$docker_enabled" == 1 ]]; then
		docker compose --project-name "$project" --file "$compose_file" down --remove-orphans --volumes || true
	fi

	if [[ -n "$oracle_mount" ]]; then
		hdiutil detach "$oracle_mount" -quiet -force || true
	fi
	rm -rf "$temporary_root"
	exit "$exit_code"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
trap 'exit 129' HUP

test_password=${GOREST_TEST_DRIVERS_PASSWORD:-GorestTest!Pass2026}
export GOREST_TEST_DRIVERS_PASSWORD="$test_password"
oracle_lib_dir=''

run_integration_test() {
	local driver=$1
	shift

	printf '\n=== Running %s migration, seeder, CRUD, and transaction checks ===\n' "$driver"
	env GOREST_DRIVER_INTEGRATION="$driver" "$@" \
		go test -v -count=1 -timeout=20m ./test/drivers -run '^TestDatabaseDriverIntegration$'
}

run_sqlite_test() {
	run_integration_test sqlite
	printf '\n=== Running SQLite route and lifecycle scenarios (excluding MongoDB and ScyllaDB tests) ===\n'
	env -u GOREST_MONGO_E2E -u GOREST_SCYLLA_E2E go test -v -count=1 -timeout=10m \
		-skip 'Native(CQL|MongoQuery)E2E' ./test/e2e
}

run_driver_test() {
	local driver=$1
	local service=$2
	local container_port=$3
	local username=$4
	local password=$5
	local database=${6:-gorest_driver_it}
	local published_port
	local extra_lib_dir=${7:-}
	local wait_timeout=360
	if [[ "$driver" == oracle ]]; then
		wait_timeout=900
	fi

	current_service="$service"
	printf '\n=== Starting disposable %s service ===\n' "$driver"
	docker compose --project-name "$project" --file "$compose_file" up \
		--detach --wait --wait-timeout "$wait_timeout" "$service"
	published_port=$(docker compose --project-name "$project" --file "$compose_file" port "$service" "$container_port")
	published_port=${published_port##*:}

	if [[ -n "$extra_lib_dir" ]]; then
		run_integration_test "$driver" \
			GOREST_DRIVER_HOST=127.0.0.1 \
			GOREST_DRIVER_PORT="$published_port" \
			GOREST_DRIVER_DATABASE="$database" \
			GOREST_DRIVER_USERNAME="$username" \
			GOREST_DRIVER_PASSWORD="$password" \
			GOREST_DRIVER_LIB_DIR="$extra_lib_dir"
	else
		run_integration_test "$driver" \
		GOREST_DRIVER_HOST=127.0.0.1 \
		GOREST_DRIVER_PORT="$published_port" \
		GOREST_DRIVER_DATABASE="$database" \
		GOREST_DRIVER_USERNAME="$username" \
		GOREST_DRIVER_PASSWORD="$password"
	fi
	case "$driver" in
		postgres|mysql|oracle|sqlserver)
			printf '\n=== Running the complete %s route and lifecycle suite ===\n' "$driver"
			GOREST_E2E_DRIVER="$driver" \
			GOREST_DRIVER_HOST=127.0.0.1 \
			GOREST_DRIVER_PORT="$published_port" \
			GOREST_DRIVER_DATABASE="$database" \
			GOREST_DRIVER_USERNAME="$username" \
			GOREST_DRIVER_PASSWORD="$password" \
			GOREST_DRIVER_LIB_DIR="$extra_lib_dir" \
			go test -v -count=1 -timeout=20m -skip 'Native(CQL|MongoQuery)E2E' ./test/e2e
			;;
		mongo|scylla)
			printf '\n=== Running the shared ORM route and lifecycle suite for %s (native SQL and NoSQL routes run in separate stages) ===\n' "$driver"
			printf 'Omitting SQL and native-driver routes (run in their own stages) plus ORM transaction-atomicity cases, which do not apply to MongoDB/ScyllaDB.\n'
			skip_pattern='Native(SQL|CQL|MongoQuery)E2E|CreateUsersInBulkRollsBackOnDatabaseErrorE2E|CreateUserWithTransactionE2E|DeleteUserWithTransactionE2E|UpdateUserWithTransactionE2E'
			if [[ "$driver" == scylla ]]; then
				printf 'Omitting ORM operations not supported by this ScyllaDB schema/API: email-keyed upsert, duplicate-key rejection, joins, non-key DISTINCT/GROUP BY, HAVING, and global ORDER BY/OFFSET.\n'
				skip_pattern+='|UpsertUserE2E|TestCreateUserE2E|GetUsersWithDepartmentE2E|GetUsersWithDistinctDepartmentIDsE2E|GetUserCountsByDepartmentE2E|GetDepartmentsByUserCountE2E|GetUsersOrderedByIDE2E|GetUsersByLimitOffsetE2E'
			fi
			GOREST_E2E_DRIVER="$driver" \
			GOREST_DRIVER_HOST=127.0.0.1 \
			GOREST_DRIVER_PORT="$published_port" \
			GOREST_DRIVER_DATABASE="$database" \
			GOREST_DRIVER_USERNAME="$username" \
			GOREST_DRIVER_PASSWORD="$password" \
			go test -v -count=1 -timeout=20m -skip "$skip_pattern" ./test/e2e
			;;
	esac

	case "$driver" in
		mongo)
			printf '\n=== Running native MongoDB route scenarios ===\n'
			GOREST_MONGO_E2E=1 GOREST_MONGO_HOST=127.0.0.1 GOREST_MONGO_PORT="$published_port" \
				GOREST_MONGO_USERNAME="$username" GOREST_MONGO_PASSWORD="$password" \
				go test -v -count=1 -timeout=10m ./test/e2e -run Mongo
			;;
		scylla)
			printf '\n=== Running native ScyllaDB CQL route scenarios ===\n'
			GOREST_SCYLLA_E2E=1 GOREST_SCYLLA_HOST=127.0.0.1 GOREST_SCYLLA_PORT="$published_port" \
				GOREST_SCYLLA_USERNAME='' GOREST_SCYLLA_PASSWORD='' \
				go test -v -count=1 -timeout=10m ./test/e2e -run CQL
			;;
	esac

	docker compose --project-name "$project" --file "$compose_file" rm --stop --force "$service"
	current_service=''
}

verify_download_checksum() {
	local archive=$1
	local expected=$2
	local actual

	if command -v sha256sum >/dev/null 2>&1; then
		actual=$(sha256sum "$archive" | awk '{print $1}')
	elif command -v shasum >/dev/null 2>&1; then
		actual=$(shasum -a 256 "$archive" | awk '{print $1}')
	else
		printf 'ERROR: sha256sum or shasum is required to verify Oracle Instant Client.\n' >&2
		return 1
	fi
	if [[ "$actual" != "$expected" ]]; then
		printf 'ERROR: Oracle Instant Client checksum mismatch (expected %s, got %s).\n' "$expected" "$actual" >&2
		return 1
	fi
}

download_oracle_client() {
	local os arch filename url checksum archive install_dir mount_dir
	os=$(uname -s)
	arch=$(uname -m)

	case "$os/$arch" in
		Linux/x86_64)
			filename='instantclient-basic-linux.x64-23.26.4.0.0.zip'
			url="https://download.oracle.com/otn_software/linux/instantclient/2326400/$filename"
			checksum='cab59aa3b411b3b50898e456fa9aabfe48d67e25cbc1bb23303800390ce7de26'
			install_dir="$temporary_root/instantclient_23_26"
			mkdir -p "$install_dir"
			archive="$temporary_root/$filename"
			printf 'Downloading Oracle Instant Client 23.26 for Linux x86_64.\n' >&2
			curl --fail --location --silent --show-error "$url" --output "$archive"
			verify_download_checksum "$archive" "$checksum"
			command -v unzip >/dev/null 2>&1 || {
				printf 'ERROR: unzip is required to prepare Oracle Instant Client.\n' >&2
				return 1
			}
			unzip -q "$archive" -d "$install_dir"
			oracle_lib_dir="$install_dir/instantclient_23_26"
			;;
		Linux/aarch64|Linux/arm64)
			filename='instantclient-basic-linux.arm64-23.26.4.0.0.zip'
			url="https://download.oracle.com/otn_software/linux/instantclient/2326400/$filename"
			checksum='688df426cabb4e5945e8409ae022acd16b92ff06582c8d93af2f41e3fb9a42f3'
			install_dir="$temporary_root/instantclient_23_26"
			mkdir -p "$install_dir"
			archive="$temporary_root/$filename"
			printf 'Downloading Oracle Instant Client 23.26 for Linux ARM64.\n' >&2
			curl --fail --location --silent --show-error "$url" --output "$archive"
			verify_download_checksum "$archive" "$checksum"
			command -v unzip >/dev/null 2>&1 || {
				printf 'ERROR: unzip is required to prepare Oracle Instant Client.\n' >&2
				return 1
			}
			unzip -q "$archive" -d "$install_dir"
			oracle_lib_dir="$install_dir/instantclient_23_26"
			;;
		Darwin/x86_64)
			filename='instantclient-basic-macos.x64-19.16.0.0.0dbru.dmg'
			url="https://download.oracle.com/otn_software/mac/instantclient/1916000/$filename"
			checksum='725eeb7c00041e926b57c3db508ea7439a1935a37f480c0ac6d46eddc4205331'
			install_dir="$temporary_root/instantclient_19_16"
			mount_dir="$temporary_root/oracle-dmg"
			mkdir -p "$mount_dir" "$install_dir"
			archive="$temporary_root/$filename"
			printf 'Downloading Oracle Instant Client 19.16 for macOS Intel.\n' >&2
			curl --fail --location --silent --show-error "$url" --output "$archive"
			verify_download_checksum "$archive" "$checksum"
			hdiutil attach -readonly -nobrowse -quiet -mountpoint "$mount_dir" "$archive"
			oracle_mount="$mount_dir"
			ditto "$mount_dir" "$install_dir"
			hdiutil detach "$oracle_mount" -quiet
			oracle_mount=''
			oracle_lib_dir="$install_dir"
			;;
		Darwin/arm64)
			filename='instantclient-basic-macos.arm64-23.26.2.0.0.dmg'
			url="https://download.oracle.com/otn_software/mac/instantclient/2326200/$filename"
			checksum='c3c4fce37a557192322717c1a8422fe098fc829a41f382f57a288b24ad6ba11e'
			install_dir="$temporary_root/instantclient_23_26"
			mount_dir="$temporary_root/oracle-dmg"
			mkdir -p "$mount_dir" "$install_dir"
			archive="$temporary_root/$filename"
			printf 'Downloading Oracle Instant Client 23.26 for macOS ARM64.\n' >&2
			curl --fail --location --silent --show-error "$url" --output "$archive"
			verify_download_checksum "$archive" "$checksum"
			hdiutil attach -readonly -nobrowse -quiet -mountpoint "$mount_dir" "$archive"
			oracle_mount="$mount_dir"
			ditto "$mount_dir" "$install_dir"
			hdiutil detach "$oracle_mount" -quiet
			oracle_mount=''
			oracle_lib_dir="$install_dir"
			;;
		*)
			printf 'ERROR: Oracle Instant Client download is not configured for %s/%s. Set GOREST_DRIVER_LIB_DIR to an existing compatible installation.\n' "$os" "$arch" >&2
			return 1
			;;
	esac
}

run_oracle_test() {
	local lib_dir=${GOREST_DRIVER_LIB_DIR:-}
	local username=${GOREST_TEST_DRIVERS_ORACLE_USERNAME:-gorest_test}
	if [[ -z "$lib_dir" ]]; then
		download_oracle_client
		lib_dir=$oracle_lib_dir
	fi
	if [[ ! -d "$lib_dir" ]]; then
		printf 'ERROR: Oracle Instant Client directory does not exist: %s\n' "$lib_dir" >&2
		return 1
	fi

	printf '\n=== Running Oracle integration checks against the disposable Oracle Free service ===\n'
	run_driver_test oracle oracle 1521 "$username" "$test_password" FREEPDB1 \
		"$lib_dir"
}

run_selected_driver() {
	case "$1" in
		sqlite) run_sqlite_test ;;
		postgres) run_driver_test postgres postgres 5432 gorest_test "$test_password" ;;
		mysql) run_driver_test mysql mysql 3306 gorest_test "$test_password" ;;
		sqlserver) run_driver_test sqlserver sqlserver 1433 sa "$test_password" gorest_driver_it ;;
		mongo) run_driver_test mongo mongo 27017 gorest_test "$test_password" ;;
		scylla) run_driver_test scylla scylla 9042 '' '' gorest_driver_it ;;
		oracle) run_oracle_test ;;
		*) printf 'ERROR: unsupported driver %q. Choose sqlite, postgres, mysql, oracle, sqlserver, mongo, or scylla.\n' "$1" >&2; exit 2 ;;
	esac
}

if [[ "$selected_driver" != all ]]; then
	run_selected_driver "$selected_driver"
	printf '\nAll %s database E2E scenarios passed.\n' "$selected_driver"
	exit 0
fi

run_selected_driver sqlite
run_selected_driver postgres
run_selected_driver mysql
run_selected_driver sqlserver
run_selected_driver mongo
run_selected_driver scylla
case "${GOREST_TEST_DRIVERS_SKIP_ORACLE:-0}" in
	0) run_selected_driver oracle; printf '\nAll seven database E2E configurations passed.\n' ;;
	1) printf '\nSix database E2E configurations passed; Oracle was skipped.\n' ;;
	*) printf 'ERROR: GOREST_TEST_DRIVERS_SKIP_ORACLE must be 0 or 1.\n' >&2; exit 1 ;;
esac
