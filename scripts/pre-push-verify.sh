#!/usr/bin/env sh

set -eu

if [ "${SKIP_VERIFY:-0}" = "1" ]; then
	echo "Skipping local verification because SKIP_VERIFY=1"
	exit 0
fi

repo_root=$(git rev-parse --show-toplevel)
cd "$repo_root"

null_sha="0000000000000000000000000000000000000000"
empty_tree=$(git hash-object -t tree /dev/null)
changed_files=""

collect_changed_files() {
	local_sha="$1"
	remote_sha="$2"

	if [ "$local_sha" = "$null_sha" ]; then
		return 0
	fi

	if [ "$remote_sha" = "$null_sha" ]; then
		base=$(git merge-base "$local_sha" origin/main 2>/dev/null || git rev-parse "${local_sha}^" 2>/dev/null || printf '%s' "$empty_tree")
	else
		base="$remote_sha"
	fi

	diff_output=$(git diff --name-only "$base" "$local_sha" || true)
	if [ -n "$diff_output" ]; then
		changed_files=$(printf '%s\n%s\n' "$changed_files" "$diff_output")
	fi
}

if [ ! -t 0 ]; then
	while read -r local_ref local_sha remote_ref remote_sha; do
		[ -z "${local_ref:-}" ] && continue
		collect_changed_files "$local_sha" "$remote_sha"
	done
fi

if [ -z "$changed_files" ]; then
	base_ref=$(git rev-parse --abbrev-ref --symbolic-full-name @{upstream} 2>/dev/null || true)
	if [ -z "$base_ref" ] && git show-ref --verify --quiet refs/remotes/origin/main; then
		base_ref="origin/main"
	fi

	if [ -n "$base_ref" ]; then
		changed_files=$(git diff --name-only "$base_ref"...HEAD || true)
	else
		changed_files=$(git diff --name-only "$empty_tree" HEAD || true)
	fi
	fi

changed_files=$(printf '%s\n' "$changed_files" | sed '/^$/d' | sort -u)

if [ -z "$changed_files" ]; then
	echo "No local changes require verification."
	exit 0
fi

run_full=0
run_go=0
run_web=0

if printf '%s\n' "$changed_files" | grep -Eq '^(Makefile|\.githooks/|web/\.husky/|scripts/pre-push-verify\.sh|\.golangci\.yml)$'; then
	run_full=1
fi

if printf '%s\n' "$changed_files" | grep -Eq '(^web/|^web/\.husky/)'; then
	run_web=1
fi

if printf '%s\n' "$changed_files" | grep -Eq '(^app/|^audit/|^blob/|^cache/|^cli/|^client/|^cmd/|^contextutil/|^crypto/|^encrypt/|^errors/|^events/|^git/|^http/|^infraprovider/|^job/|^langstats/|^livelog/|^lock/|^logging/|^mcp/|^profiler/|^pubsub/|^registry/|^resources/|^secret/|^ssh/|^store/|^stream/|^tests/|^types/|^version/|\.go$|^go\.(mod|sum|tool\.mod|tool\.sum)$|^\.golangci\.yml$)'; then
	run_go=1
fi

if [ "$run_full" -eq 1 ] || { [ "$run_go" -eq 1 ] && [ "$run_web" -eq 1 ]; }; then
	echo "Running full local verification"
	exec make verify
fi

if [ "$run_go" -eq 1 ]; then
	echo "Running Go local verification"
	exec make verify-go
fi

if [ "$run_web" -eq 1 ]; then
	echo "Running web local verification"
	exec make verify-web
fi

echo "Only non-code changes detected; skipping local verification."
