#!/usr/bin/env bash

set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
unset PI_CODING_AGENT_DIR AGY_HOME CODEX_HOME COMMITTER_PATH
fixture_root=$(mktemp -d)
export HOME="$fixture_root/home"
mkdir -p "$HOME"
trap 'rm -rf "$fixture_root"' EXIT
cp "$repo_root/AGENTS.md" "$fixture_root/agents-expected"

assert_link() {
  local path=$1
  local target=$2
  [ -L "$path" ] && [ "$(readlink "$path")" = "$target" ] || {
    printf 'unexpected symlink at %s (expected %s)\n' "$path" "$target" >&2
    exit 1
  }
}

config_for() {
  case "$1" in
    pi) printf '%s' "$2/.pi/agent" ;;
    agy) printf '%s' "$2/.gemini/antigravity-cli" ;;
    codex) printf '%s' "$2/.codex" ;;
  esac
}

for harness in pi agy codex; do
  harness_home="$fixture_root/$harness"
  config_dir=$(config_for "$harness" "$harness_home")
  HOME="$harness_home" "$repo_root/scripts/setup-$harness"
  assert_link "$config_dir/AGENTS.md" "$repo_root/instructions/global.md"
  assert_link "$config_dir/skills" "$repo_root/skills"
  assert_link "$harness_home/.local/bin/committer" "$repo_root/scripts/committer"

  HOME="$harness_home" "$repo_root/scripts/setup-harness" "$harness"
  assert_link "$config_dir/AGENTS.md" "$repo_root/instructions/global.md"
  assert_link "$config_dir/skills" "$repo_root/skills"

  custom_dir="$fixture_root/custom-$harness"
  case "$harness" in
    pi) export PI_CODING_AGENT_DIR="$custom_dir" ;;
    agy) export AGY_HOME="$custom_dir" ;;
    codex) export CODEX_HOME="$custom_dir" ;;
  esac
  HOME="$harness_home" "$repo_root/scripts/setup-harness" "$harness"
  assert_link "$custom_dir/AGENTS.md" "$repo_root/instructions/global.md"
  assert_link "$custom_dir/skills" "$repo_root/skills"

  legacy_dir="$fixture_root/legacy-$harness"
  legacy_config=$(config_for "$harness" "$legacy_dir")
  mkdir -p "$legacy_config"
  ln -s "$repo_root/AGENTS.md" "$legacy_config/AGENTS.md"
  case "$harness" in
    pi) export PI_CODING_AGENT_DIR="$legacy_config" ;;
    agy) export AGY_HOME="$legacy_config" ;;
    codex) export CODEX_HOME="$legacy_config" ;;
  esac
  HOME="$harness_home" "$repo_root/scripts/setup-harness" "$harness"
  assert_link "$legacy_config/AGENTS.md" "$repo_root/instructions/global.md"
  cmp "$repo_root/AGENTS.md" "$fixture_root/agents-expected"
done

expect_refusal() {
  local config_dir=$1
  if PI_CODING_AGENT_DIR="$config_dir" "$repo_root/scripts/setup-harness" pi >/dev/null 2>&1; then
    printf 'expected setup to refuse existing path: %s/AGENTS.md\n' "$config_dir" >&2
    exit 1
  fi
}

file_config="$fixture_root/file/.pi/agent"
mkdir -p "$file_config"
printf 'local instructions\n' > "$file_config/AGENTS.md"
cp "$file_config/AGENTS.md" "$fixture_root/file-expected"
expect_refusal "$file_config"
cmp "$file_config/AGENTS.md" "$fixture_root/file-expected"

foreign_config="$fixture_root/foreign/.pi/agent"
mkdir -p "$foreign_config" "$fixture_root/foreign-target"
printf 'do not modify this target\n' > "$fixture_root/foreign-target/sentinel"
cp "$fixture_root/foreign-target/sentinel" "$fixture_root/foreign-expected"
ln -s "$fixture_root/foreign-target/sentinel" "$foreign_config/AGENTS.md"
expect_refusal "$foreign_config"
[ "$(readlink "$foreign_config/AGENTS.md")" = "$fixture_root/foreign-target/sentinel" ]
cmp "$fixture_root/foreign-target/sentinel" "$fixture_root/foreign-expected"

dangling_config="$fixture_root/dangling/.pi/agent"
mkdir -p "$dangling_config"
ln -s "$fixture_root/missing-target" "$dangling_config/AGENTS.md"
expect_refusal "$dangling_config"
[ "$(readlink "$dangling_config/AGENTS.md")" = "$fixture_root/missing-target" ]

directory_config="$fixture_root/directory/.pi/agent"
mkdir -p "$directory_config/AGENTS.md"
printf 'directory sentinel\n' > "$directory_config/AGENTS.md/sentinel"
cp "$directory_config/AGENTS.md/sentinel" "$fixture_root/directory-expected"
expect_refusal "$directory_config"
cmp "$directory_config/AGENTS.md/sentinel" "$fixture_root/directory-expected"
