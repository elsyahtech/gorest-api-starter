#!/usr/bin/env bash

set -Eeuo pipefail

credential_fields='Password|Secret|Token|APIKey|ClientSecret|AccessKey|PrivateKey'
backup_dir=''
backup_failed=0
declare -a backup_files=()
declare -a backup_paths=()

confirm() {
	local answer
	printf '%s [y/N] ' "$1"
	read -r answer
	[[ "$answer" =~ ^[Yy]$ ]]
}

restore_originals() {
	local index
	local failed=0

	for index in "${!backup_files[@]}"; do
		if ! cp -p "${backup_paths[$index]}" "${backup_files[$index]}"; then
			printf 'ERROR: Could not restore %s from its temporary backup.\n' "${backup_files[$index]}" >&2
			failed=1
		fi
	done

	if (( failed == 0 )); then
		backup_files=()
		backup_paths=()
	else
		backup_failed=1
		return 1
	fi
}

cleanup() {
	local exit_code=$?
	trap - EXIT

	if ((${#backup_files[@]})); then
		printf 'Restoring original configuration files before exiting...\n' >&2
		if ! restore_originals; then
			exit_code=1
		fi
	fi

	if [[ -n "$backup_dir" && -d "$backup_dir" ]]; then
		if (( backup_failed == 0 )); then
			rm -rf "$backup_dir"
		else
			printf 'Temporary backups were kept at: %s\n' "$backup_dir" >&2
		fi
	fi

	exit "$exit_code"
}

unstage_changes() {
	if ! git reset -q -- .; then
		printf 'ERROR: Could not safely unstage the changes. Inspect the Git index manually.\n' >&2
		return 1
	fi
}

repo_root=$(git rev-parse --show-toplevel 2>/dev/null) || {
	printf 'ERROR: Run this command from inside a Git repository.\n' >&2
	exit 1
}
cd "$repo_root"

branch=$(git symbolic-ref --quiet --short HEAD) || {
	printf 'ERROR: Cannot push from a detached HEAD. Check out a branch first.\n' >&2
	exit 1
}

git remote get-url origin >/dev/null 2>&1 || {
	printf 'ERROR: Git remote "origin" is not configured.\n' >&2
	exit 1
}

if ! git diff --cached --quiet; then
	printf 'ERROR: The Git index already contains staged changes. Commit or unstage them before running make push.\n' >&2
	exit 1
fi

if [[ ! -d app/config ]]; then
	printf 'ERROR: Expected configuration directory app/config was not found.\n' >&2
	exit 1
fi

config_files=(app/config/*.go)
if [[ ! -f "${config_files[0]}" ]]; then
	printf 'ERROR: No Go configuration files were found in app/config.\n' >&2
	exit 1
fi

if ! command -v perl >/dev/null 2>&1; then
	printf 'ERROR: Perl is required to sanitize Go string literals.\n' >&2
	exit 1
fi

if ! confirm 'This will temporarily blank non-empty credential fields in app/config/*.go before staging. Continue?'; then
	printf 'Cancelled. No files were changed or staged.\n'
	exit 0
fi

backup_dir=$(mktemp -d "${TMPDIR:-/tmp}/gorest-push.XXXXXX")
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
trap 'exit 129' HUP

sanitized_files=0
for file in "${config_files[@]}"; do
	[[ -f "$file" && ! -L "$file" ]] || continue

	backup_path="$backup_dir/${#backup_files[@]}"
	cp -p "$file" "$backup_path"
	backup_files+=("$file")
	backup_paths+=("$backup_path")

	fields=$(
		perl -i -pe '
			s/^([ \t]*(Password|Secret|Token|APIKey|ClientSecret|AccessKey|PrivateKey)[ \t]*:[ \t]*)"((?:\\.|[^"\\])*)"/
				length($3) ? do { print STDERR "$2\n"; $1 . q{""} } : $&
			/gem;
			s/^([ \t]*(Password|Secret|Token|APIKey|ClientSecret|AccessKey|PrivateKey)[ \t]*:[ \t]*)`([^`]*)`/
				length($3) ? do { print STDERR "$2\n"; $1 . q{""} } : $&
			/gem;
		' "$file" 2>&1 >/dev/null
	)

	if [[ -n "$fields" ]]; then
		fields=${fields//$'\n'/,}
		printf 'Sanitized %s: %s\n' "$file" "${fields%,}"
		sanitized_files=$((sanitized_files + 1))
	fi
done

printf 'Sanitization complete in %s configuration file(s). Original files will be restored automatically after staging.\n' "$sanitized_files"
if ! confirm 'Do you want to stage all changes with git add -A?'; then
	printf 'Cancelled before staging. Original configuration files will be restored.\n'
	exit 0
fi

if ! git add -A; then
	printf 'ERROR: Staging failed. Attempting to unstage any partial changes.\n' >&2
	unstage_changes || true
	exit 1
fi

printf 'Restoring original configuration files; the staged versions remain sanitized.\n'
if ! restore_originals; then
	unstage_changes || true
	exit 1
fi

set +e
git grep --cached --quiet -I -E "^[[:space:]]*(${credential_fields}):[[:space:]]*(\"[^\"]+\"|\`[^\`]+\`)" -- app/config
scan_status=$?
set -e

if (( scan_status == 0 )); then
	printf 'ERROR: A non-empty credential field remains staged in app/config. Changes were not committed.\n' >&2
	unstage_changes || true
	exit 1
elif (( scan_status != 1 )); then
	printf 'ERROR: Could not verify staged configuration files. Changes were not committed.\n' >&2
	unstage_changes || true
	exit 1
fi

if ! git diff --cached --check; then
	printf 'ERROR: The staged changes did not pass git diff --check.\n' >&2
	unstage_changes || true
	exit 1
fi

if git diff --cached --quiet; then
	printf 'There are no changes to commit.\n'
	exit 0
fi

printf '\nStaged change summary:\n'
git diff --cached --stat

if ! confirm 'Do you want to commit these staged changes?'; then
	unstage_changes
	printf 'Cancelled. Staged changes were removed from the index; working files were preserved.\n'
	exit 0
fi

printf 'Commit message: '
read -r commit_message
if [[ -z "$commit_message" ]]; then
	printf 'ERROR: The commit message cannot be empty.\n' >&2
	unstage_changes
	exit 1
fi

if ! git commit -m "$commit_message"; then
	printf 'ERROR: Commit failed. Attempting to unstage the changes.\n' >&2
	unstage_changes || true
	exit 1
fi

printf 'Commit created on branch %s.\n' "$branch"
if ! confirm "Do you want to push this commit to origin/$branch?"; then
	printf 'The commit remains local and has not been pushed.\n'
	exit 0
fi

git push -u origin "$branch"
