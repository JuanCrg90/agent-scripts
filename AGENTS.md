# Agent Scripts

Repository-local guidance for this project. Shared operating policy is maintained
in [`instructions/global.md`](instructions/global.md) and linked to harnesses by
`setup-harness`.

## Layout

- `instructions/global.md`: shared agent policy.
- `skills/`: reusable skills, including `pi-delegation` for bounded standalone
  Pi workers.
- `scripts/setup-harness`: safe, idempotent harness resource setup.
- `tests/setup-harness.sh`: setup regression tests.

## Verification

```sh
pnpm test
bash -n scripts/setup-harness tests/setup-harness.sh
```

Setup links instructions and skills; it does not configure credentials, models,
MCP servers, or other machine-local settings. It migrates only the exact legacy
managed `AGENTS.md` symlink and refuses to replace files or other links.
