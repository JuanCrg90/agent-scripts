package syncer

import "path/filepath"

type Target struct {
	Name   string
	Source string
	Dest   string
	Kind   string
}

const (
	KindFile = "file"
	KindDir  = "dir"
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
			Source: filepath.Join(base, "gemini", "settings.json"),
			Dest:   filepath.Join(opts.GeminiHome, "settings.json"),
			Kind:   KindFile,
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
