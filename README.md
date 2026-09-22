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

The setup command links `AGENTS.md` and `skills/` into
`~/.pi/agent/`, plus `scripts/committer` into `~/.local/bin/committer`.
Ensure `~/.local/bin` is on your `PATH`. It is safe to run again, but refuses
to replace an existing file or a link to another source.

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
standard environment variable:

```sh
PI_CODING_AGENT_DIR="$HOME/.config/pi-agent" ./scripts/setup-harness pi
AGY_HOME="$HOME/.config/agy" ./scripts/setup-harness agy
CODEX_HOME="$HOME/.config/codex" ./scripts/setup-harness codex
```

## Browser tools

`browser-tools` drives Chrome through the DevTools Protocol without MCP. It
requires Node.js 22.12 or newer, a Chrome or Chromium browser, and `rsync` when
using `--profile`.

```sh
pnpm install --frozen-lockfile
pnpm run browser-tools start
pnpm run browser-tools nav https://example.com
pnpm run browser-tools screenshot
```

On macOS, Google Chrome is detected from `/Applications`. On Linux, the command
checks `google-chrome`, `google-chrome-stable`, `chromium`, and
`chromium-browser`. Use `--chrome-path` when your browser lives elsewhere.

## Verification

```sh
pnpm test
pnpm run browser-tools --help
```

## Repository layout

- `AGENTS.md`: shared operating instructions.
- `skills/`: reusable Agent Skills packages.
- `scripts/setup-pi`: default Pi setup.
- `scripts/setup-agy`, `scripts/setup-codex`: opt-in setup for occasional use.
- `scripts/browser-tools.ts`: Chrome DevTools CLI.
- `tools.md`: tool catalog for agents.

`.pi/` and `.antigravitycli/` are intentionally local and ignored. Keep
credentials, MCP servers, model choices, and other machine-specific settings in
the harness that uses them.
