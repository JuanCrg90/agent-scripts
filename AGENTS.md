# AGENTS.md

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

- Use herdr for orchestration; the herdr skill is your source of truth. 
- Always use openai-codex provider for GPT models.
- For development tasks, spawn new agents with the Pi harness `pi`. Use gpt-5.6-terra  for regular tasks and gpt-5.6-sol for tasks requiring high reasoning.
- For review tasks, spawn new agents with Antigravity `agy` cli as harness with the Gemini 3.1 Pro model, unless the Juan explicitly requests another model or harness. 

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
  `./scripts/committer`, use that. See `~/Projects/agent-scripts/tools.md` for
  usage and commit-message style.
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

Read `~/Projects/agent-scripts/tools.md` for the full tool catalog if it exists.

For browser automation, prefer the local `browser-tools.ts` CLI
(`~/Projects/agent-scripts/scripts/browser-tools.ts`). It uses Chrome's
DevTools Protocol directly and does not require an MCP server, matching the
no-MCP setup of this harness.

## RTK

Shell commands: prefer `rtk <cmd>`. Verify with `rtk --version`, `rtk gain`,
`which rtk`. See `~/Projects/agent-scripts/tools.md` for full `rtk` and
`committer` usage.

Use `AGENTS.md` as source of truth. Do not depend on generated `GEMINI.md`
overrides.

### trash

- Move files to Trash: `trash …` (system command).

### gh

- GitHub CLI for PRs/CI/releases. Given issue/PR URL (or `/pull/5`): use `gh`, not web search.
- Examples: `gh issue view <url> --comments -R owner/repo`, `gh pr view <url> --comments --files -R owner/repo`.
