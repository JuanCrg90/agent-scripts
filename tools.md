# Tools Reference

CLI tools available on Peter's machines. Use these for agentic tasks.

## gh

GitHub CLI for PRs, issues, CI, releases.

**Usage**: `gh help`

When someone shares a GitHub URL, use `gh` to read it:

```bash
gh issue view <url> --comments
gh pr view <url> --comments --files
gh run list / gh run view <id>
```
