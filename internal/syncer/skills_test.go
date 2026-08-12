package syncer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildPlanLinksEverySkillRoot(t *testing.T) {
	baseDir := t.TempDir()
	writeFixture(t, filepath.Join(baseDir, "AGENTS.md"), "# agents\n")
	writeFixture(t, filepath.Join(baseDir, "gemini", "settings.base.json"), "{}\n")
	writeFixture(t, filepath.Join(baseDir, "config", "mcp", "servers.json"), `{"mcpServers": {}}`)
	writeFixture(t, filepath.Join(baseDir, "skills", "demo", "SKILL.md"), "# demo\n")
	writeFixture(t, filepath.Join(baseDir, "scripts", "commiter"), "#!/bin/sh\n")

	opts := testOptions(baseDir)
	plan, err := BuildPlan(opts)
	if err != nil {
		t.Fatalf("BuildPlan returned error: %v", err)
	}
	if err := Apply(plan, opts, func(Action) (PromptDecision, error) {
		return DecisionYes, nil
	}); err != nil {
		t.Fatalf("Apply returned error: %v", err)
	}

	for _, home := range []string{
		opts.AgentsHome,
		opts.CodexHome,
		opts.GeminiHome,
		opts.AntigravityHome,
		opts.OpenCodeHome,
		opts.PiHome,
	} {
		assertSymlink(t, filepath.Join(home, "skills"), filepath.Join(baseDir, "skills"))
	}
	for _, home := range []string{
		opts.AgentsHome,
		opts.CodexHome,
		opts.GeminiHome,
		opts.AntigravityHome,
		opts.OpenCodeHome,
		opts.PiHome,
	} {
		assertSymlink(t, filepath.Join(home, "AGENTS.md"), filepath.Join(baseDir, "AGENTS.md"))
	}
}

func TestImportSkillsCopiesNewSkillsAndSkipsSystem(t *testing.T) {
	root := t.TempDir()
	repo := filepath.Join(root, "repo")
	canonical := filepath.Join(root, "canonical")
	writeFixture(t, filepath.Join(repo, "skills", "demo", "SKILL.md"), "# demo\n")
	writeFixture(t, filepath.Join(repo, "skills", ".system", "metadata"), "private\n")

	if err := ImportSkills(repo, canonical); err != nil {
		t.Fatalf("ImportSkills returned error: %v", err)
	}
	if got := readFile(t, filepath.Join(canonical, "demo", "SKILL.md")); got != "# demo\n" {
		t.Fatalf("unexpected imported skill: %q", got)
	}
	if _, err := os.Stat(filepath.Join(canonical, ".system")); !os.IsNotExist(err) {
		t.Fatalf("expected .system to be skipped, got %v", err)
	}
}

func TestImportSkillsIsAtomicWhenLaterSkillConflicts(t *testing.T) {
	root := t.TempDir()
	repo := filepath.Join(root, "repo")
	canonical := filepath.Join(root, "canonical")
	writeFixture(t, filepath.Join(repo, "skills", "new-skill", "SKILL.md"), "# new\n")
	writeFixture(t, filepath.Join(repo, "skills", "conflict", "SKILL.md"), "# imported\n")
	writeFixture(t, filepath.Join(canonical, "conflict", "SKILL.md"), "# canonical\n")

	if err := ImportSkills(repo, canonical); err == nil {
		t.Fatal("expected skill conflict")
	}
	if _, err := os.Stat(filepath.Join(canonical, "new-skill")); !os.IsNotExist(err) {
		t.Fatalf("expected no partial import, got %v", err)
	}
}

func TestImportSkillsRejectsConflicts(t *testing.T) {
	root := t.TempDir()
	repo := filepath.Join(root, "repo")
	canonical := filepath.Join(root, "canonical")
	writeFixture(t, filepath.Join(repo, "skills", "demo", "SKILL.md"), "# imported\n")
	writeFixture(t, filepath.Join(canonical, "demo", "SKILL.md"), "# canonical\n")

	if err := ImportSkills(repo, canonical); err == nil {
		t.Fatal("expected skill conflict")
	}
	if got := readFile(t, filepath.Join(canonical, "demo", "SKILL.md")); got != "# canonical\n" {
		t.Fatalf("canonical skill overwritten: %q", got)
	}
}

func testOptions(baseDir string) Options {
	return Options{
		BaseDir:         baseDir,
		AgentsHome:      filepath.Join(baseDir, "agents-home"),
		CodexHome:       filepath.Join(baseDir, "codex-home"),
		GeminiHome:      filepath.Join(baseDir, "gemini-home"),
		AntigravityHome: filepath.Join(baseDir, "antigravity-home"),
		OpenCodeHome:    filepath.Join(baseDir, "opencode-home"),
		PiHome:          filepath.Join(baseDir, "pi-home"),
	}
}

func assertSymlink(t *testing.T, path, target string) {
	t.Helper()
	link, err := os.Readlink(path)
	if err != nil {
		t.Fatalf("Readlink %s: %v", path, err)
	}
	if link != target {
		t.Fatalf("link %s = %s, want %s", path, link, target)
	}
}
