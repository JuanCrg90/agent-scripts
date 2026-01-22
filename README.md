# Agent Skills Hub

My source of truth for using agent skills with Codex, Claude, and Gemini.

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

The `agent-sync` Go CLI keeps Codex, Gemini CLI, and Antigravity in sync with this repo.

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

### Flags

- `--base <path>` override repo base
- `--codex-home <path>` override Codex home
- `--gemini-home <path>` override Gemini home
- `--symlink` use symlinks instead of copying (always in sync, but links break if you move the repo and some tools dislike symlinks)

Example with flags:

```sh
bin/agent-sync sync --base ~/Projects/agent-scripts --codex-home ~/.codex --gemini-home ~/.gemini
```

Symlink mode:

```sh
bin/agent-sync sync --symlink
```

### What syncs

- `AGENTS.md` → `~/.codex/AGENTS.md`, `~/.gemini/AGENTS.md`, `~/.gemini/GEMINI.md`
- `gemini/settings.json` → `~/.gemini/settings.json`
- `skills/` → `~/.codex/skills`, `~/.gemini/antigravity/skills`
- `scripts/commiter` → `~/.codex/scripts/commiter`

## Acknowledgment

Inspired by @steipete agent-scripts repository <https://github.com/steipete/agent-scripts/>.
Sandi Metz rules skill reference by @nateberkopec <https://github.com/nateberkopec/dotfiles/tree/main/files/home/.claude/skills/sandi-metz-rules>.
