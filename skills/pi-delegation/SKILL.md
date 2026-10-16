---
name: pi-delegation
description: Delegate bounded work to a standalone Pi process. Use when the primary coordinator launches a Pi worker outside Herdr; Herdr is explicit opt-in and follows its own skill.
---

# Standalone Pi delegation

This skill guides the primary coordinator launching a fresh worker. The
coordinator remains responsible for task design, results, and review. Only the
spawned session with an explicit worker role is a bounded worker: it must not
delegate, coordinate, create branches/worktrees, commit, or change machine-local
configuration. It works only in assigned scope/files and reports blockers.

## Preflight and launch

Set `PROVIDER`, `MODEL`, and `THINKING` to the exact values assigned by the
coordinator; the example uses `medium`. The coordinator must verify the chosen
model supports the intended thinking level: Pi may clamp unsupported levels, so
successful output parsing alone does not establish that the requested level was
used. Do not silently substitute a preferred model if unavailable. Preflight
under the same offline, model, thinking, and resource flags as launch, then
verify the output contains a row whose provider and model ID both exactly
match. A successful exit or the presence of other models does not establish
availability.

```sh
cd /absolute/path/to/workspace
PROVIDER='openai-codex'
MODEL='exact-model-id-from-coordinator'
THINKING='medium'
TASK='/absolute/path/to/task.json'
RUN_DIR=$(mktemp -d)
chmod 700 "$RUN_DIR"
pi --offline --provider "$PROVIDER" --model "$MODEL" \
  --no-extensions --no-mcp --no-skills --no-prompt-templates --no-themes \
  --thinking "$THINKING" --list-models >"$RUN_DIR/preflight.out" 2>"$RUN_DIR/preflight.err"
```

The list command can print multiple models; inspect `preflight.out` for the
exact provider/ID match. If the requested row is absent, stop and report that;
do not reroute without coordinator authorization. If a provider deliberately
requires an extension, explicitly add that same extension to both preflight and
launch and explain the exception to the coordinator.

Use Pi's one-shot print mode, a fresh in-memory session, and the task file as
context. This read-only tool list replaces defaults and exposes only file
inspection/search tools. For a writer, set `TOOLS=read,grep,find,ls,bash,edit,write`
instead; plain `--tools` replaces defaults, so only enable write-capable tools
for assigned files. Run exactly once in the foreground and capture stdout,
stderr, and exit status:

```sh
TOOLS=read,grep,find,ls
set +e
pi --offline --provider "$PROVIDER" --model "$MODEL" \
  --no-extensions --no-mcp --no-skills --no-prompt-templates --no-themes \
  --thinking "$THINKING" --tools "$TOOLS" --print --no-session "@$TASK" \
  >"$RUN_DIR/stdout" 2>"$RUN_DIR/stderr"
status=$?
set -e
```

Retain `RUN_DIR` until results are validated. Never resume or fork the parent's session. `--offline` disables
startup network activity, not model provider traffic. Tools are not an OS
sandbox: they inherit the Pi process's permissions. `--no-approve` is a project
trust gate, not execution permission or isolation. Limit credentials and files;
use an OS/container sandbox when isolation is required.

## Task and result contract

Pass a concise JSON task via `@task.json`, including worker role, goal, workspace,
assigned files, constraints, and checks. Example:

```json
{
  "role": "spawned implementation worker; bounded task; do not delegate",
  "goal": "Implement the assigned regression fix",
  "workspace": "/absolute/path/to/repository",
  "files": ["src/example.ts", "tests/example.test.ts"],
  "constraints": ["No commits or branches", "Preserve unrelated changes"],
  "checks": ["Run the focused test"],
  "result": "Return exactly Amadeus: on line one, then one JSON object with status, files, checks, blockers, and risks. No fences or prose."
}
```

Require final output to be exactly `Amadeus:` on the first line followed by one
JSON object (no code fence or surrounding prose). Strip only that exact first
line before parsing the JSON. Require `status` to be `completed`, `blocked`, or
`failed`, plus `files`, `checks`, `blockers`, and `risks`; validate values and
reject empty, partial, or invalid results. The prefix is a human-facing response
convention, not part of the JSON. Print mode provides final text only; output
truncated at the model token limit is nonzero and must fail validation. JSON
mode is a JSONL event stream: validate completion and final response events,
not merely a zero exit. Prefer print mode for this contract.

If no suitable worker/model is available, the coordinator implements directly
and discloses unavailable independent review. Keep one writer per overlapping
file.

## Completion and recovery

Wait for process exit; retain stdout, stderr, and status. A nonzero exit, empty
output, malformed result, incomplete required fields, or non-completed status
is not success. Inspect captured output and the workspace diff before deciding
next steps. On timeout, do not blindly resubmit: inspect partial edits and
terminate only the child process you own if necessary. Never kill unrelated
processes.

## Pi references

Recipes follow the Pi 1.2.0 installed documentation: `docs/cli.md` (options),
`docs/cli-integration.md` (print and JSON completion semantics),
`docs/configuration.md` (resource discovery), and `docs/security.md` (trust and
process permissions). Check `pi --help` for the installed CLI's exact options.
