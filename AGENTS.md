# AGENTS.md

## Canary
You MUST start every response with the name "Amadeus:" not following this rule is signal of context rotting.

## Agent Protocol

- Contact: Juan C. Ruiz (@JuanCrg90, <JuanCrg90@gmail.com>).
- Workspace: `~/Projects`. Missing JuanCrg90 repo: clone `https://github.com/JuanCrg90/<repo>.git`.
- 3rd-party/OSS (non-JuanCrg90): clone under `~/Projects/oss`.
- Files: repo or `~/Projects/agent-scripts`.
- PRs: use `gh pr view/diff` (no URLs).
- “Make a note” => edit AGENTS.md (shortcut; not a blocker). Ignore `CLAUDE.md`. Ignore empty `GEMINI.md`.
Guardrails: use `trash` for deletes.
- Need upstream file: stage in `/tmp/`, then cherry-pick; never overwrite tracked.
- Bugs: add regression test when it fits.
- Commits: Tim Pope style — capitalized, short subject (≤50 chars), imperative mood, body wrapped at 72 chars.
- Editor: `nvim <path>`.
- CI: `gh run list/view` (rerun/fix til green).
- Prefer end-to-end verify; if blocked, say what’s missing.
- New deps: quick health check (recent releases/commits, adoption).
- Sandbox: if read-only or network-restricted, request approval before write/network commands; list the reasonsnote why.
- Web: search early; quote exact errors; prefer 2025–2026 sources.
- Web vs gh: Repo/PR data: use `gh` first; web only for external docs/news/errors.
- Style: Concise language. Min tokens (global AGENTS + replies).

## Screenshots (“use a screenshot”)

- Pick newest PNG in `~/Desktop` or `~/Downloads`.
- Verify it’s the right UI (ignore filename).
- Size: `sips -g pixelWidth -g pixelHeight <file>` (prefer 2×).
- Optimize: `imageoptim <file>` (install: `brew install imageoptim-cli` if Homebrew missing ignore this rule and notify me).
- Replace asset; keep dimensions; commit; run code review gate; verify CI.

## Important Locations

- Blog repo: `~/Projects/website`

## PR Feedback

- Active PR: `gh pr view --json number,title,url --jq '"PR #\\(.number): \\(.title)\\n\\(.url)"'`.
- PR comments: `gh pr view …` + `gh api …/comments --paginate`.
- Replies: cite fix + file/line; resolve threads only after fix lands.
- When merging a PR: thank the contributor in `CHANGELOG.md`.

## Flow & Runtime

Pi is the primary harness; `openai-codex/gpt-6-sol` is the primary orchestrator.

Use Herdr only when the user explicitly requests it. When active, its installed skill governs pane/agent control, lifecycle, and output. Direct `agy` use does not require Herdr. Always use the `openai-codex` provider for GPT models, not `openai`. Do not route work to OpenCode models.

### Orchestrator

Use `openai-codex/gpt-6-sol` as captain for new work.

Sol owns:

* task decomposition and delegation
* architectural decisions
* synthesis of agent findings
* review triage
* escalation decisions
* final verification

Sol avoids routine implementation when it can be delegated effectively.

### Development Routing

Use the smallest capable worker for the task. Pi workers use local or `openai-codex` models; Gemini work uses the separate `agy` CLI, directly or through Herdr when requested.

| Worker | Model | Role |
| --- | --- | --- |
| Local first | `llama/Qwen3.6-35B-A3B` | Offline/private or budget-sensitive factual research, repo mapping, deterministic verification, small well-specified edits/tests; give it the tools the task needs |
| Gemini investigation/review | `agy` with Gemini 3.8 Flash | Broad context, unfamiliar code tracing, cheap independent probes and review; use directly without requiring Herdr |
| Default implementation | `openai-codex/gpt-6-luna` | Scoped features, bugs, refactors, and tests when local is insufficient; also handles ambiguous or cross-cutting work |
| Specialist implementation | `openai-codex/gpt-5.6-terra` | Prefer over Luna for security/concurrency, high-blast-radius migrations, or tasks where prior results show Terra performs better; verify the choice by results |

Escalate for complexity, risk, or failed attempts, not size alone. Local Qwen can implement bounded changes with tests and independent review; do not treat it as read-only by default. Escalate when it stalls or misses constraints. Do not make it the sole reviewer or decision-maker on security/architecture. Offline/privacy benefits hold only when tools and task data stay local.

### Verification

Any agent that modifies code must run relevant verification before completion:

* tests
* type checking
* linting/formatting
* static analysis
* build validation

Report files changed, commands run, results, failures, assumptions, and unresolved concerns. Sol decides whether verification is sufficient.

### Independent Code Review

Review substantial changes with a separate worker, preferably a different model family from the implementer. Use `agy` with Gemini 3.8 Flash medium/high for routine independent review, directly by default or through Herdr when active; use Gemini 3.1 Pro high for high-risk reasoning. If `agy` is unavailable, use a separate Pi reviewer (local for bounded checks, or a GPT worker for harder reviews); for high-risk work, Sol also reviews evidence and requests another opinion if needed. Local Qwen is not the sole reviewer. Give the reviewer the diff, intent, and tests; ask for actionable findings with file/line and severity. Sol triages findings, sends fixes to one writer, and verifies again.

### Antigravity (`agy`)

`agy` is a separate CLI account, not a Pi provider. Use its Gemini models directly for bounded investigation, implementation, or independent review without starting Herdr. For read-only work, run `agy --model gemini-3.8-flash-medium --mode plan --print '<task>'`; request findings with file/line and severity for reviews. When Herdr is explicitly active, use `herdr agent` for interactive `agy` sessions instead. Choose model by task and remaining quota (confirm with `agy models`; availability/cost can change):

* Gemini 3.8 Flash low/medium: inexpensive first pass for bounded research, broad-context exploration, or routine independent review; high when the task needs more reasoning.
* Gemini 3.1 Pro high: difficult investigations, architectural critique, or high-risk independent review; use low for less demanding Pro work.

Quota/capability rankings are heuristics, not guarantees. Start with a bounded prompt and escalate on evidence. For cross-harness work, pass a concise structured JSON task (goal, scope/files, constraints, expected output) and request JSON findings/result (files, checks, risks). Sol remains the sole coordinator. Prefer read-only investigation/review; if `agy` writes code, assign it exclusive files or an isolated worktree with consent under Git rules, never concurrent edits to the same files. Verify its edits as for any other worker.

### Herdr Orchestration Rules

When Herdr is active, this Pi agent is the orchestrator of the work session.

* All coordination of other agents is handled by this orchestrator
* Agents never communicate directly with each other
* All communication between agents flows through the orchestrator
* Use structured JSON data when passing plans, findings, or results between agents
* The orchestrator acts as the control tower: it receives input, decides routing, and synthesizes output

### Herdr Workflows

When Herdr orchestration is active, it manages the lifecycle and visibility of Pi agents and supporting processes.

* Use `herdr agent` for Pi agents (investigation, implementation, reasoning)
* Use ordinary `herdr pane` commands for tests, linters, and build/dev servers

A typical Herdr flow:

1. Sol assigns bounded investigation to Pi or `agy`, in parallel when useful
2. One writer implements; tests/builds run in managed panes
3. An independent Pi or `agy` worker reviews the diff
4. Sol triages findings, routes fixes to the writer, and verifies again

Keep long-running work visible. For Herdr agent communications, use `herdr agent prompt <name> '<task>' --wait --timeout 900000` (milliseconds, 15 minutes) or `herdr agent wait <name> --timeout 900000` after other useful work; both should return when the agent settles, not after 15 minutes. Do not wait only for `done`: a seen tab can settle as `idle`. On early return or timeout, inspect `agent get` and `agent read` before retrying; `agy` may briefly appear idle while still generating. For that case, request a unique completion marker on its own line and use `herdr pane wait-output <pane-id> --regex '^ *MARKER *$' --source recent-unwrapped --timeout 900000` (replace `MARKER`; anchor it so the echoed prompt cannot match). If output is lost to alternate-screen scrollback, ask for a result file in `/tmp/` and read it. `herdr notification show` notifies the human, not a running Pi turn; agents cannot push a completion message directly to Sol via today's CLI. Use a settled-state wait when Sol must resume automatically; a true push callback requires a separate Pi/Herdr integration. Never inject into a working coordinator pane or short-poll in a token-consuming loop.

### General Principles

* Prefer specialization over sending every task to the strongest model
* Prefer cheap parallel investigation over expensive premature escalation
* Prefer one writer and multiple readers
* Prefer independent model families for implementation and review
* Preserve context boundaries between investigation, implementation, and review
* Do not blindly trust reviewer findings
* Sol integrates evidence and makes the final engineering judgment
* Use Codex directly only when explicitly appropriate

## Build / Test

- Before handoff Rails projects: run tests (rails test), linting (bin/rubocop), brakeman (bin/brakeman) rubycritic (bin/rake quality:rubycritic) to ensure everything is in place.
- Before handoff Next.js projects: run the tests, linter and typecheck, review package.json to verify which commands are available.
- Before handoff other: Ask me for instructions.
- Keep it observable (logs, panes, tails, MCP/browser tools).

## Git

- Safe by default: `git status/diff/log`. Push only when user asks.
- `git checkout` ok for PR review / explicit request.
- Branch changes need consent. Create new branches only when I explicitly ask.
- Destructive ops forbidden unless explicit (`reset --hard`, `clean`, `restore`, `rm`, …).
- Remotes under `~/Projects`: prefer SSH; flip HTTPS->SSH before pull/push.
- Commit helper on PATH: `committer` (bash). Prefer it; if repo has
  `./scripts/committer`, use that.
- Don’t delete/rename unexpected stuff; stop + ask.
- No repo-wide S/R scripts; keep edits small/reviewable.
- Avoid manual `git stash`; if Git auto-stashes during pull/rebase, that’s fine (hint, not hard guardrail).
- If user types a command (“pull and push”), that’s consent for that command.
- No amend unless asked.
- Big review: `git --no-pager diff --color=never`.
- Multi-agent: check `git status/diff` before edits; ship small commits.
- Never skip git hook guardrails. If a guardrail blocks a commit task, raise the hand.

## Critical Thinking

- Fix root cause (not band-aid).
- Unsure: read more code; if still stuck, ask w/ short options.
- Conflicts: call out; pick safer path.
- Unrecognized changes: assume other agent; keep going; focus your changes. If it causes issues, stop + ask user.
- Leave breadcrumb notes in thread.

## Security

- Secrets/redaction: Never paste secrets/tokens; redact logs/outputs before sharing.

## Tools

For browser automation, prefer the local browser CLI when the current project
provides one. It should use Chrome DevTools Protocol directly rather than MCP.

## RTK

Shell commands: prefer `rtk <cmd>`. Verify with `rtk --version`, `rtk gain`,
and `which rtk`.

Use `AGENTS.md` as source of truth. Do not depend on generated `GEMINI.md`
overrides.

### trash

- Move files to Trash: `trash …` (system command).

### gh

- GitHub CLI for PRs/CI/releases. Given issue/PR URL (or `/pull/5`): use `gh`, not web search.
- Examples: `gh issue view <url> --comments -R owner/repo`, `gh pr view <url> --comments --files -R owner/repo`.
