package syncer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	codexManagedBlockStart = "# agent-sync: managed mcp start"
	codexManagedBlockEnd   = "# agent-sync: managed mcp end"
)

type manifestFile struct {
	Servers map[string]MCPServer `json:"mcpServers"`
}

type MCPServer struct {
	Type    string            `json:"type,omitempty"`
	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	URL     string            `json:"url,omitempty"`
}

func RenderGeminiSettings(baseDir string) ([]byte, error) {
	basePath := filepath.Join(baseDir, "gemini", "settings.base.json")
	manifestPath := filepath.Join(baseDir, "config", "mcp", "servers.json")

	baseData, err := os.ReadFile(basePath)
	if err != nil {
		return nil, err
	}

	var settings map[string]any
	if err := json.Unmarshal(baseData, &settings); err != nil {
		return nil, fmt.Errorf("parse %s: %w", basePath, err)
	}

	manifest, err := readManifest(manifestPath)
	if err != nil {
		return nil, err
	}

	mcpServers := make(map[string]any, len(manifest.Servers))
	for name, server := range manifest.Servers {
		entry, err := server.geminiEntry()
		if err != nil {
			return nil, fmt.Errorf("render Gemini MCP server %s: %w", name, err)
		}
		mcpServers[name] = entry
	}
	settings["mcpServers"] = mcpServers

	rendered, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(rendered, '\n'), nil
}

func RenderAgyMcpConfig(baseDir string) ([]byte, error) {
	manifestPath := filepath.Join(baseDir, "config", "mcp", "servers.json")
	manifest, err := readManifest(manifestPath)
	if err != nil {
		return nil, err
	}

	mcpServers := make(map[string]any, len(manifest.Servers))
	for name, server := range manifest.Servers {
		entry := make(map[string]any)
		switch server.kind() {
		case "stdio":
			entry["command"] = server.Command
			entry["args"] = server.ArgsOrEmpty()
			if len(server.Env) > 0 {
				entry["env"] = server.Env
			}
		case "http":
			entry["serverUrl"] = server.URL
		}
		mcpServers[name] = entry
	}

	config := map[string]any{
		"mcpServers": mcpServers,
	}

	rendered, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(rendered, '\n'), nil
}

func RenderCodexConfig(baseDir string, existing []byte) ([]byte, error) {
	manifestPath := filepath.Join(baseDir, "config", "mcp", "servers.json")
	manifest, err := readManifest(manifestPath)
	if err != nil {
		return nil, err
	}

	clean := stripManagedCodexBlock(string(existing))
	if err := ensureNoCodexConflicts(clean, manifest); err != nil {
		return nil, err
	}

	block, err := renderCodexManagedBlock(manifest)
	if err != nil {
		return nil, err
	}

	clean = strings.TrimRight(clean, "\n")
	if block == "" {
		if clean == "" {
			return []byte{}, nil
		}
		return []byte(clean + "\n"), nil
	}

	if clean == "" {
		return []byte(block), nil
	}
	return []byte(clean + "\n\n" + block), nil
}

func readManifest(path string) (manifestFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return manifestFile{}, err
	}

	var manifest manifestFile
	if err := json.Unmarshal(data, &manifest); err != nil {
		return manifestFile{}, fmt.Errorf("parse %s: %w", path, err)
	}

	for name, server := range manifest.Servers {
		if err := server.validate(); err != nil {
			return manifestFile{}, fmt.Errorf("invalid MCP server %s: %w", name, err)
		}
	}

	return manifest, nil
}

func (server MCPServer) validate() error {
	serverType := server.kind()
	switch serverType {
	case "stdio":
		if server.Command == "" {
			return fmt.Errorf("stdio server requires command")
		}
	case "http":
		if server.URL == "" {
			return fmt.Errorf("http server requires url")
		}
	default:
		return fmt.Errorf("unsupported type %q", server.Type)
	}
	return nil
}

func (server MCPServer) kind() string {
	if server.Type == "" {
		return "stdio"
	}
	return server.Type
}

func (server MCPServer) geminiEntry() (map[string]any, error) {
	switch server.kind() {
	case "stdio":
		entry := map[string]any{
			"command": server.Command,
			"args":    server.ArgsOrEmpty(),
		}
		if len(server.Env) > 0 {
			entry["env"] = server.Env
		}
		return entry, nil
	case "http":
		return map[string]any{"url": server.URL}, nil
	default:
		return nil, fmt.Errorf("unsupported type %q", server.Type)
	}
}

func (server MCPServer) ArgsOrEmpty() []string {
	if server.Args == nil {
		return []string{}
	}
	return server.Args
}

func renderCodexManagedBlock(manifest manifestFile) (string, error) {
	if len(manifest.Servers) == 0 {
		return "", nil
	}

	names := sortedServerNames(manifest.Servers)
	var b strings.Builder
	b.WriteString(codexManagedBlockStart)
	b.WriteString("\n")

	for index, name := range names {
		server := manifest.Servers[name]
		if index > 0 {
			b.WriteString("\n")
		}
		switch server.kind() {
		case "stdio":
			b.WriteString(fmt.Sprintf("[mcp_servers.%s]\n", name))
			b.WriteString(fmt.Sprintf("command = %s\n", tomlString(server.Command)))
			if len(server.ArgsOrEmpty()) > 0 {
				b.WriteString(fmt.Sprintf("args = %s\n", tomlStringArray(server.ArgsOrEmpty())))
			}
			if len(server.Env) > 0 {
				b.WriteString("\n")
				b.WriteString(fmt.Sprintf("[mcp_servers.%s.env]\n", name))
				for _, key := range sortedKeys(server.Env) {
					b.WriteString(fmt.Sprintf("%s = %s\n", key, tomlString(server.Env[key])))
				}
			}
		case "http":
			b.WriteString(fmt.Sprintf("[mcp_servers.%s]\n", name))
			b.WriteString(fmt.Sprintf("url = %s\n", tomlString(server.URL)))
		default:
			return "", fmt.Errorf("unsupported type %q", server.Type)
		}
	}

	b.WriteString("\n")
	b.WriteString(codexManagedBlockEnd)
	b.WriteString("\n")
	return b.String(), nil
}

func ensureNoCodexConflicts(config string, manifest manifestFile) error {
	for _, name := range sortedServerNames(manifest.Servers) {
		header := fmt.Sprintf("[mcp_servers.%s]", name)
		envHeader := fmt.Sprintf("[mcp_servers.%s.env]", name)
		if strings.Contains(config, header) || strings.Contains(config, envHeader) {
			return fmt.Errorf("codex config already defines managed MCP server %s outside the agent-sync block", name)
		}
	}
	return nil
}

func stripManagedCodexBlock(config string) string {
	start := strings.Index(config, codexManagedBlockStart)
	if start == -1 {
		return config
	}

	end := strings.Index(config[start:], codexManagedBlockEnd)
	if end == -1 {
		return strings.TrimRight(config[:start], "\n")
	}

	end += start + len(codexManagedBlockEnd)
	before := strings.TrimRight(config[:start], "\n")
	after := strings.TrimLeft(config[end:], "\n")

	switch {
	case before == "" && after == "":
		return ""
	case before == "":
		return after
	case after == "":
		return before
	default:
		return before + "\n\n" + after
	}
}

func buildManagedAction(target Target, opts Options) (Action, error) {
	var (
		rendered []byte
		err      error
	)

	switch target.Kind {
	case KindGeminiConfig:
		rendered, err = RenderGeminiSettings(opts.BaseDir)
	case KindAgyMcpConfig:
		rendered, err = RenderAgyMcpConfig(opts.BaseDir)
	case KindCodexConfig:
		existing, readErr := os.ReadFile(target.Dest)
		if readErr != nil && !os.IsNotExist(readErr) {
			return Action{}, readErr
		}
		rendered, err = RenderCodexConfig(opts.BaseDir, existing)
	default:
		return Action{}, fmt.Errorf("unknown managed target kind: %s", target.Kind)
	}
	if err != nil {
		return Action{}, err
	}

	return buildContentAction(target.Source, target.Dest, rendered), nil
}

func buildContentAction(src, dest string, rendered []byte) Action {
	existing, err := os.ReadFile(dest)
	switch {
	case isSymlink(dest):
		return Action{Type: ActionUpdate, Source: src, Dest: dest, Content: rendered, Note: "replace symlink"}
	case err == nil && bytes.Equal(existing, rendered):
		return Action{Type: ActionNoop, Source: src, Dest: dest, Content: rendered}
	case err == nil:
		return Action{Type: ActionUpdate, Source: src, Dest: dest, Content: rendered, Note: "rendered"}
	case os.IsNotExist(err) && len(rendered) == 0:
		return Action{Type: ActionNoop, Source: src, Dest: dest, Content: rendered}
	case os.IsNotExist(err):
		return Action{Type: ActionCreate, Source: src, Dest: dest, Content: rendered, Note: "rendered"}
	default:
		return Action{Type: ActionUpdate, Source: src, Dest: dest, Content: rendered, Note: err.Error()}
	}
}

func tomlString(value string) string {
	encoded, _ := json.Marshal(value)
	return string(encoded)
}

func tomlStringArray(values []string) string {
	var encoded []string
	for _, value := range values {
		encoded = append(encoded, tomlString(value))
	}
	return "[" + strings.Join(encoded, ", ") + "]"
}

func sortedServerNames(servers map[string]MCPServer) []string {
	names := make([]string, 0, len(servers))
	for name := range servers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func sortedKeys(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
