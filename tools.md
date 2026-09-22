# Tools Reference

CLI tools available on Juan's machines. Use these for agentic tasks.

## gh

GitHub CLI for PRs, issues, CI, releases.

**Usage**: `gh help`

When someone shares a GitHub URL, use `gh` to read it:

```bash
gh issue view <url> --comments
gh pr view <url> --comments --files
gh run list / gh run view <id>
```

## browser-tools

Lightweight Chrome automation via the DevTools Protocol. Lives at
`scripts/browser-tools.ts` in this repository.

**Why this exists**: this harness intentionally avoids MCPs, and the installed
browser skills (`browser-testing-with-devtools`, `web-perf`) require an MCP
server. This tool talks directly to Chrome over CDP using `puppeteer-core`, so
it works without any MCP setup.

**Prerequisites**:

- Node.js 22.12 or newer and pnpm 12.5.1, then `pnpm install --frozen-lockfile` in this repository.
- Google Chrome on macOS, or Google Chrome/Chromium on Linux. Override detection
  with `--chrome-path`.
- `rsync` only when using `start --profile`.

**Run it** (from this repository's root):

```bash
pnpm run browser-tools <command>
```

**Commands**:

| Command | Purpose |
| --- | --- |
| `start` | Launch Chrome with remote debugging on port 9222. |
| `nav <url>` | Navigate the active tab (or `--new` for a new tab). |
| `eval <code>` | Run JavaScript in the active page. |
| `screenshot` | Capture the viewport and print the temp PNG path. |
| `pick <message>` | Interactive DOM picker; click to print element metadata. |
| `console` | Capture/tail console logs. |
| `search <query>` | Google search with optional readable content extraction. |
| `content <url>` | Extract readable article content as markdown. |
| `cookies` | Dump cookies from the active tab as JSON. |
| `inspect` | List Chrome debug instances and their tabs. |
| `kill` | Terminate Chrome debug instances. |

**Common workflow**:

```bash
# 1. Start Chrome
pnpm run browser-tools start

# 2. Navigate somewhere
pnpm run browser-tools nav https://gema.cafe

# 3. Screenshot (prints path; copy to Desktop if needed)
pnpm run browser-tools screenshot

# 4. Evaluate JS in the active page
pnpm run browser-tools eval "document.title"
```

**Notes**:

- Always operate on the **last active tab**.
- Screenshots are saved to `/tmp`; move them manually (e.g. to `~/Desktop`).
- For pages that need the user's profile (logins, cookies), start with
  `--profile`.
- To shut down the debug Chrome instance: `pnpm run browser-tools kill --all`.

## rtk

Preferred wrapper for shell commands in the agent workflow.

- Shell commands: prefer `rtk <cmd>`.
- Examples:
  - `rtk git status`
  - `rtk cargo test`
  - `rtk npm run build`
  - `rtk pytest -q`
- Meta commands:
  - `rtk gain`
  - `rtk gain --history`
  - `rtk proxy <cmd>`
- Verify install: `rtk --version`, `rtk gain`, `which rtk`.

## committer

Commit helper on PATH. Stages only listed paths; required here. A repo may
also ship its own `./scripts/committer`.

**Usage**:

```bash
committer -m "subject" [--dry-run] [--force] "file" ["file" ...]
```

- `-m` takes multi-line messages: first line = subject, rest = body.
- `--dry-run`: show staged diff without committing.
- `--force`: remove stale git lock.
- `.` is disallowed; list specific file paths only.

**Commit message style** (Tim Pope):

- Subject line: capitalized, imperative mood, ≤50 chars.
  - ✅ `Add login page`
  - ❌ `Added login page`, ❌ `feat: add login`
- Blank line separating subject from body (unless body is omitted).
- Body wrapped at 72 chars. Explain *what* and *why*, not *how*.

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
```
