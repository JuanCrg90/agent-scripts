# AGENTS.md

## Agent Protocol

- Contact: Juan C. Ruiz (@JuanCrg90, <JuanCrg90@gmail.com>).
- Workspace: `~/Projects`. Missing JuanCrg90 repo: clone `https://github.com/JuanCrg90/<repo>.git`.
- 3rd-party/OSS (non-JuanCrg90): clone under `~/Projects/oss`.
- Files: repo or `~/Projects/agent-scripts`.
- PRs: use `gh pr view/diff` (no URLs).
- “Make a note” => edit AGENTS.md (shortcut; not a blocker). Ignore `CLAUDE.md`. Ignore empty `GEMINI.md`.
- No `./runner`. Guardrails: use `trash` for deletes.
- Need upstream file: stage in `/tmp/`, then cherry-pick; never overwrite tracked.
- Bugs: add regression test when it fits.
- Keep files <~500 LOC; split/refactor as needed.
- Commits: Tim Pope style — capitalized, short subject (≤50 chars), imperative mood, body wrapped at 72 chars.
- Editor: `nvim <path>`.
- CI: `gh run list/view` (rerun/fix til green).
- Prefer end-to-end verify; if blocked, say what’s missing.
- New deps: quick health check (recent releases/commits, adoption).
- Sandbox: if read-only or network-restricted, request approval before write/network commands; note why.
- Web: search early; quote exact errors; prefer 2025–2026 sources.
- Web vs gh: Repo/PR data: use `gh` first; web only for external docs/news/errors.
- Style: telegraph. Drop filler/grammar. Min tokens (global AGENTS + replies).

## Screenshots (“use a screenshot”)

- Pick newest PNG in `~/Desktop` or `~/Downloads`.
- Verify it’s the right UI (ignore filename).
- Size: `sips -g pixelWidth -g pixelHeight <file>` (prefer 2×).
- Optimize: `imageoptim <file>` (install: `brew install imageoptim-cli` if Homebrew missing ignore this rule and notify me).
- Replace asset; keep dimensions; commit; run gate; verify CI.

## Important Locations

- Blog repo: `~/Projects/website`

## PR Feedback

- Active PR: `gh pr view --json number,title,url --jq '"PR #\\(.number): \\(.title)\\n\\(.url)"'`.
- PR comments: `gh pr view …` + `gh api …/comments --paginate`.
- Replies: cite fix + file/line; resolve threads only after fix lands.
- When merging a PR: thank the contributor in `CHANGELOG.md`.

## Flow & Runtime

- Use herdr for orchestation, spin new pi agents with gpt-5.6 luna for regular tasks, use gpt-5.6 sol for tasks that requires high reasoning

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
- Commit helper on PATH: `committer` (bash). Prefer it; if repo has `./scripts/committer`, use that.
- Don’t delete/rename unexpected stuff; stop + ask.
- No repo-wide S/R scripts; keep edits small/reviewable.
- Avoid manual `git stash`; if Git auto-stashes during pull/rebase, that’s fine (hint, not hard guardrail).
- If user types a command (“pull and push”), that’s consent for that command.
- No amend unless asked.
- Big review: `git --no-pager diff --color=never`.
- Multi-agent: check `git status/diff` before edits; ship small commits.

## Language/Stack Notes

- Ruby On Rails: For Rails projects use rails-mcp-server to follow Rails community best practices, if it's not installed skip it

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

## RTK

- RTK in workflow. Shell commands: prefer `rtk <cmd>`.
- Examples: `rtk git status`, `rtk cargo test`, `rtk npm run build`, `rtk pytest -q`.
- Meta: `rtk gain`, `rtk gain --history`, `rtk proxy <cmd>`.
- Verify install: `rtk --version`, `rtk gain`, `which rtk`.
- Use `AGENTS.md` as source of truth. Do not depend on generated `GEMINI.md` overrides.

### committer

- Commit helper (PATH). Stages only listed paths; required here. Repo may also ship `./scripts/committer`.
- Usage: `committer -m "subject" [--dry-run] [--force] "file" ["file" ...]`
- `-m` takes multi-line messages: first line = subject, rest = body.
- Flags: `--dry-run` (show staged diff), `--force` (remove stale git lock).
- `.` is disallowed; list specific file paths only.

#### Commit message style (Tim Pope)

- **Subject line:** capitalized, imperative mood, ≤50 chars.
  - ✅ `Add login page`
  - ❌ `Added login page`, ❌ `feat: add login`
- **Blank line** separating subject from body (unless body is omitted).
- **Body:** wrapped at 72 chars. Explain *what* and *why*, not *how*.
  - More detailed explanatory text, if necessary.
  - Further paragraphs come after blank lines.
  - Bullet points are okay, too:
    - Use a hyphen or asterisk, followed by a single space.
    - Use a hanging indent.

Examples:

```
Add login page

Add a login page with JWT auth and refresh token rotation.

- Uses the new auth service from commit abc123
- Adds routes to config/router.ex
```

```
Fix deadlock in database pool

The connection pool was deadlocking because checkout was called
without a timeout, blocking when all connections were in use.

Wrap checkout in a timeout and add telemetry to detect future
contention.

### trash

- Move files to Trash: `trash …` (system command).

### gh

- GitHub CLI for PRs/CI/releases. Given issue/PR URL (or `/pull/5`): use `gh`, not web search.
- Examples: `gh issue view <url> --comments -R owner/repo`, `gh pr view <url> --comments --files -R owner/repo`.
