# Agent Scripts

A Pi-first, versioned source of truth for agent instructions, skills, and local
CLI tools. Pi discovers the linked `AGENTS.md` and `skills/` directory natively;
this repository deliberately does not manage shared MCP configuration.

## Quick start

Clone the repository anywhere, then link it to Pi:

```sh
git clone git@github.com:JuanCrg90/agent-scripts.git
cd agent-scripts
./scripts/setup-pi
```

The setup command links `instructions/global.md` as `AGENTS.md` and `skills/`
into `~/.pi/agent/`, plus `scripts/committer` into
`~/.local/bin/committer`. Ensure `~/.local/bin` is on your `PATH`. It is safe
to rerun. Existing installations should rerun `./scripts/setup-pi` to migrate
before reloading. Only the exact legacy managed symlink from this repository's
old `AGENTS.md` is migrated; files, directories, and other (including dangling)
links are preserved and cause setup to stop.

## Optional harnesses

Pi is the default. Configure the occasional harness only on machines where you
use it:

```sh
./scripts/setup-agy
./scripts/setup-codex
```

These commands link the same instructions and skills into the harness's normal
configuration directory. They do not configure models, authentication, MCPs,
or other machine-local settings.

For a custom config directory, use the shared command with the harness's
standard environment variable. After changing linked instructions or skills,
use `/reload` in Pi or restart the process to load updated resources:

```sh
PI_CODING_AGENT_DIR="$HOME/.config/pi-agent" ./scripts/setup-harness pi
AGY_HOME="$HOME/.config/agy" ./scripts/setup-harness agy
CODEX_HOME="$HOME/.config/codex" ./scripts/setup-harness codex
```

## Browser automation

Prefer `playwright-cli` for browser manipulation. It does not require MCP
configuration. Install it separately using the
[Playwright CLI skill](skills/playwright-cli/SKILL.md), which also
covers browser setup, interactions, snapshots, and testing.

```sh
playwright-cli open https://example.com
playwright-cli snapshot
playwright-cli screenshot
playwright-cli close
```

Chrome DevTools MCP remains optional for specialized diagnostics that Playwright
cannot cover, when configured.

## Verification

```sh
pnpm test
```

## Repository layout

- `AGENTS.md`: short repository-local map and checks.
- `instructions/global.md`: shared operating instructions linked to harnesses.
- `skills/`: reusable Agent Skills packages; `skills/pi-delegation/` documents
  bounded standalone Pi workers (Herdr is explicit opt-in).
- `scripts/setup-pi`: default Pi setup.
- `scripts/setup-agy`, `scripts/setup-codex`: opt-in setup for occasional use.
- `tools.md`: tool catalog for agents.

`.pi/` and `.antigravitycli/` are intentionally local and ignored. Keep
credentials, MCP servers, model choices, and other machine-specific settings in
the harness that uses them.
