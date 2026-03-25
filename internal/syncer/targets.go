package syncer

import "path/filepath"

type Target struct {
	Name    string
	Source  string
	Sources []string
	Dest    string
	Kind    string
}

const (
	KindFile         = "file"
	KindDir          = "dir"
	KindGeminiConfig = "gemini-config"
	KindCodexConfig  = "codex-config"
)

func Targets(opts Options) []Target {
	base := opts.BaseDir
	return []Target{
		{
			Name:   "agents-codex",
			Source: filepath.Join(base, "AGENTS.md"),
			Dest:   filepath.Join(opts.CodexHome, "AGENTS.md"),
			Kind:   KindFile,
		},
		{
			Name:   "agents-gemini",
			Source: filepath.Join(base, "AGENTS.md"),
			Dest:   filepath.Join(opts.GeminiHome, "AGENTS.md"),
			Kind:   KindFile,
		},
		{
			Name:   "gemini-context",
			Source: filepath.Join(base, "AGENTS.md"),
			Dest:   filepath.Join(opts.GeminiHome, "GEMINI.md"),
			Kind:   KindFile,
		},
		{
			Name:   "gemini-settings",
			Source: "gemini/settings.base.json + config/mcp/servers.json",
			Sources: []string{
				filepath.Join(base, "gemini", "settings.base.json"),
				filepath.Join(base, "config", "mcp", "servers.json"),
			},
			Dest: filepath.Join(opts.GeminiHome, "settings.json"),
			Kind: KindGeminiConfig,
		},
		{
			Name:   "codex-mcp",
			Source: "config/mcp/servers.json",
			Sources: []string{
				filepath.Join(base, "config", "mcp", "servers.json"),
			},
			Dest: filepath.Join(opts.CodexHome, "config.toml"),
			Kind: KindCodexConfig,
		},
		{
			Name:   "skills-codex",
			Source: filepath.Join(base, "skills"),
			Dest:   filepath.Join(opts.CodexHome, "skills"),
			Kind:   KindDir,
		},
		{
			Name:   "skills-antigravity",
			Source: filepath.Join(base, "skills"),
			Dest:   filepath.Join(opts.GeminiHome, "antigravity", "skills"),
			Kind:   KindDir,
		},
		{
			Name:   "commiter",
			Source: filepath.Join(base, "scripts", "commiter"),
			Dest:   filepath.Join(opts.CodexHome, "scripts", "commiter"),
			Kind:   KindFile,
		},
	}
}
