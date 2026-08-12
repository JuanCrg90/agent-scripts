# Agent Skills Hub

My source of truth for using agent skills with Codex, Claude, Gemini, Antigravity, OpenCode, and Pi.
MCP server definitions live here once and sync to supported harness configs.

## Purpose

- Single place for agent workflows, scripts, and reference docs.
- Keep instructions consistent across tools and models.

## Repo Structure

- `AGENTS.md`: top-level operating instructions.
- `tools.md`: tool catalog and usage notes.
- `docs/`: guides, conventions, and references.
- `scripts/`: helper scripts for agent tasks.
- `skills/`: reusable skill definitions and references.

## Usage

- Start with `AGENTS.md`.
- Use `tools.md` for tooling and command expectations.
- Add new docs under `docs/` and scripts under `scripts/`.

## Sync CLI

The `agent-sync` Go CLI keeps Codex, Gemini CLI, Antigravity, OpenCode, and Pi in sync with this repo.

### Build

```sh
go build -o bin/agent-sync ./cmd/agent-sync
```

### Commands

```sh
bin/agent-sync status
bin/agent-sync diff
bin/agent-sync plan
bin/agent-sync doctor
bin/agent-sync init
bin/agent-sync sync
```

Command meaning:

- `status`: summary counts (create/update/delete/noop)
- `diff`: list planned changes (dry-run)
- `plan`: alias for `diff`
- `doctor`: check missing source files; non-zero exit if missing
- `init`: create target dirs then run sync
- `sync`: apply changes; prompts on overwrite/delete

### Defaults

- Base: `~/Projects/agent-scripts`
- Codex home: `$CODEX_HOME` or `~/.codex`
- Gemini home: `~/.gemini`
- Antigravity home: `~/.gemini/antigravity-cli`
- OpenCode config dir: `~/.config/opencode`
- Pi agent config dir: `~/.pi/agent`

### Flags

- `--base <path>` override repo base
- `--codex-home <path>` override Codex home
- `--gemini-home <path>` override Gemini home
- `--antigravity-home <path>` override Antigravity CLI home
- `--opencode-home <path>` override OpenCode config dir
- `--pi-home <path>` override Pi agent config dir
- `--agents-home <path>` override general Agents config dir

Example with flags:

```sh
bin/agent-sync sync --base ~/Projects/agent-scripts --codex-home ~/.codex --gemini-home ~/.gemini --opencode-home ~/.config/opencode --pi-home ~/.pi/agent
```

Skills and `AGENTS.md` are always directory/file symlinks, respectively.
This makes `agent-scripts` the only source of truth. Do not run
`npx skills add`; use the native import command instead:

```sh
bin/agent-sync add https://github.com/cloudflare/skills
```

`add` clones the source into a temporary directory, imports its `skills/*`
folders into `agent-scripts/skills`, skips harness `.system` metadata, rejects
same-name content conflicts, then links every supported skill root.

### What syncs

- `AGENTS.md` → `~/.agents/AGENTS.md`, `~/.codex/AGENTS.md`, `~/.gemini/AGENTS.md`, `~/.gemini/antigravity-cli/AGENTS.md`, `~/.config/opencode/AGENTS.md`, `~/.pi/agent/AGENTS.md`
- `gemini/settings.base.json` + `config/mcp/servers.json` → rendered `~/.gemini/settings.json`
- `config/mcp/servers.json` → managed MCP block in `~/.codex/config.toml`
- `AGENTS.md` → `~/.config/opencode/AGENTS.md`
- `config/mcp/servers.json` → merged `mcp` entries in `~/.config/opencode/opencode.json`
- `skills/` → `~/.agents/skills`, `~/.codex/skills`, `~/.gemini/skills`, `~/.gemini/antigravity-cli/skills`, `~/.config/opencode/skills`, `~/.pi/agent/skills`
- `scripts/commiter` → `~/.codex/scripts/commiter`

### RTK

RTK guidance now lives in repo `AGENTS.md`, so `agent-sync` keeps it in both Codex and Gemini without writing `~/.gemini/GEMINI.md`.
That avoids deleting RTK-managed Gemini overrides while keeping `AGENTS.md` as the shared source of truth.

### Shared MCP workflow

Edit shared servers in `config/mcp/servers.json`.
Keep Gemini-only non-MCP settings in `gemini/settings.base.json`.
Remote MCP servers can include `"httpHeaders"` for committed static headers and `"httpHeadersEnv"` for local-only secrets loaded from environment variables during sync.

Then review and apply:

```sh
bin/agent-sync plan
bin/agent-sync sync
```

Verify parity after sync:

```sh
codex mcp list
jq '.mcpServers' ~/.gemini/settings.json
```

Codex MCP config is merged into a managed block inside `~/.codex/config.toml`.
Leave other Codex settings outside that block so local model/trust config stays intact.

OpenCode MCP config is merged into `~/.config/opencode/opencode.json` under `mcp`.
Existing non-MCP settings and unrelated MCP server names are preserved; server names from `config/mcp/servers.json` are refreshed from the shared manifest.

## Acknowledgment

Inspired by @steipete agent-scripts repository <https://github.com/steipete/agent-scripts/>.
Sandi Metz rules skill reference by @nateberkopec <https://github.com/nateberkopec/dotfiles/tree/main/files/home/.claude/skills/sandi-metz-rules>.
