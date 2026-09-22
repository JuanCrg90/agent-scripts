#!/usr/bin/env bash

set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
temp_dir=$(mktemp -d)
trap 'rm -rf "$temp_dir"' EXIT

assert_link() {
  local path=$1
  local target=$2

  [ -L "$path" ] || {
    printf 'expected symlink: %s\n' "$path" >&2
    exit 1
  }
  [ "$(readlink "$path")" = "$target" ] || {
    printf 'unexpected target for %s: %s\n' "$path" "$(readlink "$path")" >&2
    exit 1
  }
}

for harness in pi agy codex; do
  home="$temp_dir/$harness"
  HOME="$home" "$repo_root/scripts/setup-$harness"

  case "$harness" in
    pi) config_dir="$home/.pi/agent" ;;
    agy) config_dir="$home/.gemini/antigravity-cli" ;;
    codex) config_dir="$home/.codex" ;;
  esac

  assert_link "$config_dir/AGENTS.md" "$repo_root/AGENTS.md"
  assert_link "$config_dir/skills" "$repo_root/skills"
  assert_link "$home/.local/bin/committer" "$repo_root/scripts/committer"
done

custom_pi_home="$temp_dir/custom-pi"
PI_CODING_AGENT_DIR="$custom_pi_home" "$repo_root/scripts/setup-pi"
assert_link "$custom_pi_home/AGENTS.md" "$repo_root/AGENTS.md"
assert_link "$custom_pi_home/skills" "$repo_root/skills"

existing_home="$temp_dir/existing"
mkdir -p "$existing_home/.pi/agent"
printf 'local instructions\n' > "$existing_home/.pi/agent/AGENTS.md"
if HOME="$existing_home" "$repo_root/scripts/setup-harness" pi >/dev/null 2>&1; then
  printf 'expected setup to refuse replacing a local file\n' >&2
  exit 1
fi
