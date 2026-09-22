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

Pi is the primary harness and GPT-5.6 Sol is the primary orchestrator.

Use herdr only when the user explicitly requests orchestration with herdr; when herdr is active, its installed skill is the source of truth for commands, pane management, agent lifecycle, and output collection.

Always use the `openai-codex` provider for GPT models.

### Orchestrator

Use `openai-codex/gpt-5.6-sol` as the primary orchestrator/captain.

Sol owns:

* task decomposition and delegation
* architectural decisions
* synthesis of agent findings
* review triage
* escalation decisions
* final verification

Sol avoids routine implementation when it can be delegated effectively.

### Development Routing

Use the smallest capable Pi worker for the task.

All OpenCode Go models below are accessed through Pi via the installed `pi-opencode-bridge` package, with IDs under the `oc-sdk-go` provider:

| Worker | Model | Role |
| --- | --- | --- |
| Reconnaissance | `oc-sdk-go/deepseek-v4.1-flash` | Locate files, trace flow, read-only investigation, cheap parallel probes |
| Broad context | `oc-sdk-go/qwen3.8-flash` | Understand unfamiliar subsystems, frontend/UI-heavy analysis, second investigation path |
| Default implementation | `oc-sdk-go/kimi-k2.7-code` | Well-scoped features, bug fixes, refactors, tests, routine backend/frontend work |
| Complex implementation | `openai-codex/gpt-5.6-terra` | Cross-cutting changes, high risk, security/concurrency, architectural migrations |
| Difficult escalation | `oc-sdk-go/glm-5.3` | Hard debugging or stalled implementation paths |

Escalate to Terra or GLM because of reasoning complexity, ambiguity, blast radius, security/concurrency risk, or failed attempts; do not escalate merely because a task is large.

### Verification

Any agent that modifies code must run relevant verification before completion:

* tests
* type checking
* linting/formatting
* static analysis
* build validation

Report files changed, commands run, results, failures, assumptions, and unresolved concerns. Sol decides whether verification is sufficient.

### Code Review

Open Code Review (`ocr`) is the default independent reviewer for substantial changes. It runs as a managed process, not as a Pi coding agent.

OCR is configured to use the OpenCode Go subscription:

* Provider: `opencode-go`
* Endpoint: `https://opencode.ai/zen/go/v1`
* Default review model: `qwen3.8-flash`
* High-risk review model: `qwen3.8-max`

When using OCR from Pi, the model IDs above map to Pi's `oc-sdk-go/qwen3.8-flash` and `oc-sdk-go/qwen3.8-max` via the `pi-opencode-bridge` package.

Run it with agent-oriented output and a concise background:

```bash
ocr review --audience agent --format json --background "<review context>"
```

For large output, write to a file and return its path. OCR findings are evidence; Sol triages and decides which are actionable.

### Independent Review and Gemini

Use Antigravity `agy` with Gemini 3.1 Pro only as an additional independent reviewer or investigator, not as part of the normal implementation/review path. Use it when:

* OCR and the implementation worker disagree on a significant issue
* Sol wants a second independent opinion on a high-risk finding
* the user explicitly requests Gemini or Antigravity

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
* Use ordinary `herdr pane` commands for tests, linters, build/dev servers, and OCR

A typical Herdr flow:

1. Sol acts as captain
2. Investigation workers gather context in parallel when useful
3. One implementation worker performs the change
4. Verification runs in managed panes
5. OCR runs in a review pane against the resulting diff
6. Sol collects findings and routes accepted fixes back to the implementation worker
7. Verification runs again
8. Sol decides whether another review pass is needed

Keep long-running or parallel work visible in separate panes.

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
