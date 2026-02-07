#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

CACHE_BASE="${MLIB_GO_TEST_CACHE_BASE:-/tmp/mlib_go_t4_go_test}"
GO_TEST_TIMEOUT="${GO_TEST_TIMEOUT:-20m}"
BONE_PACKAGE="${MLIB_BONE_PACKAGE:-./pkg/domain/deform}"

usage() {
	cat <<'USAGE'
WSL向け Go テスト実行ヘルパー

使い方:
  run_go_test_wsl.sh changed [go test args...]
  run_go_test_wsl.sh pkg <package> [go test args...]
  run_go_test_wsl.sh bone [go test args...]
  run_go_test_wsl.sh all [go test args...]
  run_go_test_wsl.sh all-split [go test args...]
  run_go_test_wsl.sh list-changed
  run_go_test_wsl.sh clean-testcache
  run_go_test_wsl.sh help

環境変数:
  MLIB_GO_TEST_CACHE_BASE  キャッシュ配置先 (default: /tmp/mlib_go_t4_go_test)
  GO_TEST_TIMEOUT          go test -timeout の値 (default: 20m)
  MLIB_BONE_PACKAGE        bone モード対象 package (default: ./pkg/domain/deform)
USAGE
}

setup_env() {
	mkdir -p \
		"${CACHE_BASE}/gopath" \
		"${CACHE_BASE}/gocache" \
		"${CACHE_BASE}/gomodcache" \
		"${CACHE_BASE}/gotmp"

	export GOPATH="${CACHE_BASE}/gopath"
	export GOCACHE="${CACHE_BASE}/gocache"
	export GOMODCACHE="${CACHE_BASE}/gomodcache"
	export GOTMPDIR="${CACHE_BASE}/gotmp"
}

log_env() {
	echo "[wsl-go-test] repo=${REPO_ROOT}"
	echo "[wsl-go-test] GOPATH=${GOPATH}"
	echo "[wsl-go-test] GOCACHE=${GOCACHE}"
	echo "[wsl-go-test] GOMODCACHE=${GOMODCACHE}"
	echo "[wsl-go-test] GOTMPDIR=${GOTMPDIR}"
	echo "[wsl-go-test] timeout=${GO_TEST_TIMEOUT}"
}

run_go_test() {
	local pkg="$1"
	shift || true
	local cmd=(go test -count=1 -timeout "${GO_TEST_TIMEOUT}" "${pkg}")
	if [[ "$#" -gt 0 ]]; then
		cmd+=("$@")
	fi
	echo "[wsl-go-test] running: ${cmd[*]}"
	(
		cd "${REPO_ROOT}"
		"${cmd[@]}"
	)
}

list_changed_packages() {
	(
		cd "${REPO_ROOT}"
		{
			git diff --name-only -- '*.go'
			git diff --name-only --cached -- '*.go'
			git ls-files --others --exclude-standard -- '*.go'
		} | sed '/^$/d' | while IFS= read -r file; do
			dir="$(dirname "${file}")"
			if [[ -d "${dir}" ]]; then
				printf './%s\n' "${dir}"
			fi
		done | sort -u
	)
}

run_changed_packages() {
	mapfile -t pkgs < <(list_changed_packages)
	if [[ "${#pkgs[@]}" -eq 0 ]]; then
		echo "[wsl-go-test] 変更された Go ファイルがないため、テストを実行しません。"
		return 0
	fi
	for pkg in "${pkgs[@]}"; do
		run_go_test "${pkg}" "$@"
	done
}

run_all_split() {
	mapfile -t pkgs < <(
		cd "${REPO_ROOT}"
		go list ./...
	)
	for pkg in "${pkgs[@]}"; do
		run_go_test "${pkg}" "$@"
	done
}

main() {
	local mode="${1:-help}"
	local pkg=""
	shift || true

	setup_env
	log_env

	case "${mode}" in
	changed)
		run_changed_packages "$@"
		;;
	pkg)
		if [[ "$#" -lt 1 ]]; then
			echo "[wsl-go-test] pkg モードは package 指定が必要です。"
			usage
			exit 1
		fi
		pkg="$1"
		shift || true
		run_go_test "${pkg}" "$@"
		;;
	bone)
		run_go_test "${BONE_PACKAGE}" "$@"
		;;
	all)
		run_go_test ./... "$@"
		;;
	all-split)
		run_all_split "$@"
		;;
	list-changed)
		list_changed_packages
		;;
	clean-testcache)
		(
			cd "${REPO_ROOT}"
			go clean -testcache
		)
		echo "[wsl-go-test] go clean -testcache を実行しました。"
		;;
	help|-h|--help)
		usage
		;;
	*)
		echo "[wsl-go-test] 不明なモード: ${mode}"
		usage
		exit 1
		;;
	esac
}

main "$@"
