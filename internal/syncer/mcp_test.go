package syncer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderGeminiSettings(t *testing.T) {
	baseDir := t.TempDir()
	writeFixture(t, filepath.Join(baseDir, "gemini", "settings.base.json"), `{
  "context": {"fileName": "AGENTS.md"},
  "general": {"previewFeatures": true}
}`)
	writeFixture(t, filepath.Join(baseDir, "config", "mcp", "servers.json"), `{
  "mcpServers": {
    "echo": {
      "type": "stdio",
      "command": "/bin/echo",
      "args": ["one"],
      "env": {"FOO": "bar"}
    }
  }
}`)

	rendered, err := RenderGeminiSettings(baseDir)
	if err != nil {
		t.Fatalf("RenderGeminiSettings returned error: %v", err)
	}

	var settings map[string]any
	if err := json.Unmarshal(rendered, &settings); err != nil {
		t.Fatalf("json.Unmarshal returned error: %v", err)
	}

	if settings["context"] == nil {
		t.Fatalf("expected base settings to be preserved")
	}

	mcpServers, ok := settings["mcpServers"].(map[string]any)
	if !ok {
		t.Fatalf("expected mcpServers map, got %T", settings["mcpServers"])
	}

	echo, ok := mcpServers["echo"].(map[string]any)
	if !ok {
		t.Fatalf("expected echo server entry, got %T", mcpServers["echo"])
	}

	if echo["command"] != "/bin/echo" {
		t.Fatalf("expected command to match manifest, got %v", echo["command"])
	}
}

func TestRenderCodexConfigReplacesManagedBlock(t *testing.T) {
	baseDir := t.TempDir()
	writeFixture(t, filepath.Join(baseDir, "config", "mcp", "servers.json"), `{
  "mcpServers": {
    "echo": {
      "type": "stdio",
      "command": "/bin/echo",
      "args": ["one", "two"],
      "env": {"FOO": "bar"}
    }
  }
}`)

	existing := strings.Join([]string{
		`model = "gpt-5.4"`,
		"",
		codexManagedBlockStart,
		"[mcp_servers.old]",
		`command = "/bin/old"`,
		codexManagedBlockEnd,
		"",
	}, "\n")

	rendered, err := RenderCodexConfig(baseDir, []byte(existing))
	if err != nil {
		t.Fatalf("RenderCodexConfig returned error: %v", err)
	}

	output := string(rendered)
	if !strings.Contains(output, `model = "gpt-5.4"`) {
		t.Fatalf("expected unmanaged config to remain: %s", output)
	}
	if strings.Contains(output, "[mcp_servers.old]") {
		t.Fatalf("expected old managed block to be replaced: %s", output)
	}
	if !strings.Contains(output, "[mcp_servers.echo]") {
		t.Fatalf("expected rendered managed server: %s", output)
	}
	if !strings.Contains(output, `[mcp_servers.echo.env]`) {
		t.Fatalf("expected env table for managed server: %s", output)
	}
}

func TestRenderCodexConfigRejectsUnmanagedConflicts(t *testing.T) {
	baseDir := t.TempDir()
	writeFixture(t, filepath.Join(baseDir, "config", "mcp", "servers.json"), `{
  "mcpServers": {
    "echo": {
      "type": "stdio",
      "command": "/bin/echo"
    }
  }
}`)

	_, err := RenderCodexConfig(baseDir, []byte("[mcp_servers.echo]\ncommand = \"/bin/manual\"\n"))
	if err == nil {
		t.Fatalf("expected unmanaged conflict error")
	}
}

func TestBuildPlanAndApplyManagedConfigs(t *testing.T) {
	baseDir := t.TempDir()
	codexHome := filepath.Join(baseDir, "codex-home")
	geminiHome := filepath.Join(baseDir, "gemini-home")

	writeFixture(t, filepath.Join(baseDir, "AGENTS.md"), "# agents\n## RTK\n- RTK in workflow.\n")
	writeFixture(t, filepath.Join(baseDir, "gemini", "settings.base.json"), `{
  "context": {"fileName": "AGENTS.md"}
}`)
	writeFixture(t, filepath.Join(baseDir, "config", "mcp", "servers.json"), `{
  "mcpServers": {
    "railsMcpServer": {
      "type": "stdio",
      "command": "/bin/echo",
      "args": ["rails"]
    }
  }
}`)
	writeFixture(t, filepath.Join(baseDir, "skills", "demo", "SKILL.md"), "# skill\n")
	writeFixture(t, filepath.Join(baseDir, "scripts", "commiter"), "#!/bin/sh\n")
	writeFixture(t, filepath.Join(codexHome, "config.toml"), "model = \"gpt-5.4\"\n")

	opts := Options{
		BaseDir:    baseDir,
		CodexHome:  codexHome,
		GeminiHome: geminiHome,
	}

	plan, err := BuildPlan(opts)
	if err != nil {
		t.Fatalf("BuildPlan returned error: %v", err)
	}

	if err := Apply(plan, opts, func(Action) (PromptDecision, error) {
		return DecisionYes, nil
	}); err != nil {
		t.Fatalf("Apply returned error: %v", err)
	}

	geminiSettings := readFile(t, filepath.Join(geminiHome, "settings.json"))
	if !strings.Contains(geminiSettings, `"mcpServers"`) {
		t.Fatalf("expected rendered Gemini MCP config: %s", geminiSettings)
	}

	codexConfig := readFile(t, filepath.Join(codexHome, "config.toml"))
	if !strings.Contains(codexConfig, "[mcp_servers.railsMcpServer]") {
		t.Fatalf("expected rendered Codex MCP config: %s", codexConfig)
	}
	if !strings.Contains(codexConfig, `model = "gpt-5.4"`) {
		t.Fatalf("expected existing Codex settings to remain: %s", codexConfig)
	}

	geminiAgents := readFile(t, filepath.Join(geminiHome, "AGENTS.md"))
	if !strings.Contains(geminiAgents, "RTK in workflow") {
		t.Fatalf("expected synced AGENTS.md to carry RTK guidance: %s", geminiAgents)
	}
}

func TestBuildPlanReplacesSymlinkedManagedFile(t *testing.T) {
	baseDir := t.TempDir()
	codexHome := filepath.Join(baseDir, "codex-home")
	geminiHome := filepath.Join(baseDir, "gemini-home")

	writeFixture(t, filepath.Join(baseDir, "AGENTS.md"), "# agents\n")
	writeFixture(t, filepath.Join(baseDir, "gemini", "settings.base.json"), `{
  "context": {"fileName": "AGENTS.md"}
}`)
	writeFixture(t, filepath.Join(baseDir, "config", "mcp", "servers.json"), `{
  "mcpServers": {
    "railsMcpServer": {
      "type": "stdio",
      "command": "/bin/echo"
    }
  }
}`)
	writeFixture(t, filepath.Join(baseDir, "skills", "demo", "SKILL.md"), "# skill\n")
	writeFixture(t, filepath.Join(baseDir, "scripts", "commiter"), "#!/bin/sh\n")

	if err := os.MkdirAll(geminiHome, 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	linkTarget := filepath.Join(baseDir, "generated-settings.json")
	writeFixture(t, linkTarget, "{}\n")
	if err := os.Symlink(linkTarget, filepath.Join(geminiHome, "settings.json")); err != nil {
		t.Fatalf("Symlink returned error: %v", err)
	}

	opts := Options{
		BaseDir:    baseDir,
		CodexHome:  codexHome,
		GeminiHome: geminiHome,
	}

	plan, err := BuildPlan(opts)
	if err != nil {
		t.Fatalf("BuildPlan returned error: %v", err)
	}

	found := false
	for _, action := range plan.Actions {
		if action.Dest == filepath.Join(geminiHome, "settings.json") {
			found = action.Note == "replace symlink"
		}
	}
	if !found {
		t.Fatalf("expected managed Gemini settings action to replace symlink")
	}

	if err := Apply(plan, opts, func(Action) (PromptDecision, error) {
		return DecisionYes, nil
	}); err != nil {
		t.Fatalf("Apply returned error: %v", err)
	}

	info, err := os.Lstat(filepath.Join(geminiHome, "settings.json"))
	if err != nil {
		t.Fatalf("Lstat returned error: %v", err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		t.Fatalf("expected settings.json symlink to be replaced with regular file")
	}
}

func TestBuildPlanLeavesExistingGeminiOverrideAlone(t *testing.T) {
	baseDir := t.TempDir()
	codexHome := filepath.Join(baseDir, "codex-home")
	geminiHome := filepath.Join(baseDir, "gemini-home")

	writeFixture(t, filepath.Join(baseDir, "AGENTS.md"), "# agents\n## RTK\n- RTK in workflow.\n")
	writeFixture(t, filepath.Join(baseDir, "gemini", "settings.base.json"), `{
  "context": {"fileName": "AGENTS.md"}
}`)
	writeFixture(t, filepath.Join(baseDir, "config", "mcp", "servers.json"), `{
  "mcpServers": {}
}`)
	writeFixture(t, filepath.Join(baseDir, "skills", "demo", "SKILL.md"), "# skill\n")
	writeFixture(t, filepath.Join(baseDir, "scripts", "commiter"), "#!/bin/sh\n")
	writeFixture(t, filepath.Join(geminiHome, "GEMINI.md"), "# external rtk override\n")

	opts := Options{
		BaseDir:    baseDir,
		CodexHome:  codexHome,
		GeminiHome: geminiHome,
	}

	plan, err := BuildPlan(opts)
	if err != nil {
		t.Fatalf("BuildPlan returned error: %v", err)
	}

	for _, action := range plan.Actions {
		if action.Dest == filepath.Join(geminiHome, "GEMINI.md") {
			t.Fatalf("expected GEMINI.md override to be unmanaged, got action: %+v", action)
		}
	}

	if err := Apply(plan, opts, func(Action) (PromptDecision, error) {
		return DecisionYes, nil
	}); err != nil {
		t.Fatalf("Apply returned error: %v", err)
	}

	override := readFile(t, filepath.Join(geminiHome, "GEMINI.md"))
	if override != "# external rtk override\n" {
		t.Fatalf("expected GEMINI.md override to remain untouched, got: %s", override)
	}
}

func writeFixture(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	return string(data)
}
