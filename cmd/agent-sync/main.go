package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/JuanCrg90/agent-scripts/internal/syncer"
)

const baseDefault = "~/Projects/agent-scripts"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	command := os.Args[1]
	if command == "-h" || command == "--help" || command == "help" {
		usage()
		return
	}

	fs := flag.NewFlagSet(command, flag.ExitOnError)
	base := fs.String("base", baseDefault, "Base directory for agent-scripts")
	agentsHome := fs.String("agents-home", "~/.agents", "General agents home directory")
	codexHome := fs.String("codex-home", defaultCodexHome(), "Codex home directory")
	geminiHome := fs.String("gemini-home", "~/.gemini", "Gemini home directory")
	antigravityHome := fs.String("antigravity-home", "~/.gemini/antigravity-cli", "Antigravity CLI home directory")
	opencodeHome := fs.String("opencode-home", "~/.config/opencode", "OpenCode config directory")
	piHome := fs.String("pi-home", "~/.pi/agent", "Pi agent config directory")
	fs.Parse(os.Args[2:])

	opts := syncer.Options{
		BaseDir:         expandPath(*base),
		AgentsHome:      expandPath(*agentsHome),
		CodexHome:       expandPath(*codexHome),
		GeminiHome:      expandPath(*geminiHome),
		AntigravityHome: expandPath(*antigravityHome),
		OpenCodeHome:    expandPath(*opencodeHome),
		PiHome:          expandPath(*piHome),
	}

	switch command {
	case "status":
		runStatus(opts)
	case "diff", "plan":
		runDiff(opts)
	case "doctor":
		runDoctor(opts)
	case "init":
		runInit(opts)
	case "sync":
		runSync(opts)
	case "add":
		runAdd(opts, fs.Args())
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Println("agent-sync <command> [flags]")
	fmt.Println("")
	fmt.Println("Commands: status, diff, plan, doctor, init, sync, add")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  --base <path>        Base directory (default: ~/Projects/agent-scripts)")
	fmt.Println("  --agents-home <path> General agents home (default: ~/.agents)")
	fmt.Println("  --codex-home <path>  Codex home (default: $CODEX_HOME or ~/.codex)")
	fmt.Println("  --gemini-home <path> Gemini home (default: ~/.gemini)")
	fmt.Println("  --antigravity-home <path> Antigravity CLI home (default: ~/.gemini/antigravity-cli)")
	fmt.Println("  --opencode-home <path> OpenCode config dir (default: ~/.config/opencode)")
	fmt.Println("  --pi-home <path>     Pi agent config dir (default: ~/.pi/agent)")
	fmt.Println("")
	fmt.Println("add <git-url> imports skills into canonical storage, then links every harness.")
}

func runStatus(opts syncer.Options) {
	plan, err := syncer.BuildPlan(opts)
	if err != nil {
		fatal(err)
	}
	fmt.Println(syncer.FormatSummary(syncer.Summarize(plan)))
}

func runDiff(opts syncer.Options) {
	plan, err := syncer.BuildPlan(opts)
	if err != nil {
		fatal(err)
	}
	for _, line := range syncer.FormatActions(plan) {
		fmt.Println(line)
	}
}

func runDoctor(opts syncer.Options) {
	missing := findMissingSources(opts)
	if len(missing) == 0 {
		fmt.Println("ok")
		return
	}
	for _, path := range missing {
		fmt.Printf("missing: %s\n", path)
	}
	os.Exit(1)
}

func runInit(opts syncer.Options) {
	if err := os.MkdirAll(opts.CodexHome, 0o755); err != nil {
		fatal(err)
	}
	if err := os.MkdirAll(opts.GeminiHome, 0o755); err != nil {
		fatal(err)
	}
	if err := os.MkdirAll(opts.AntigravityHome, 0o755); err != nil {
		fatal(err)
	}
	if err := os.MkdirAll(opts.OpenCodeHome, 0o755); err != nil {
		fatal(err)
	}
	if err := os.MkdirAll(opts.PiHome, 0o755); err != nil {
		fatal(err)
	}
	runSync(opts)
}

func runAdd(opts syncer.Options, args []string) {
	if len(args) != 1 {
		fatal(fmt.Errorf("add requires exactly one Git repository URL"))
	}

	dir, err := os.MkdirTemp("", "agent-sync-add-*")
	if err != nil {
		fatal(err)
	}
	defer os.RemoveAll(dir)

	sourceDir := filepath.Join(dir, "source")
	if err := exec.Command("git", "clone", "--depth", "1", args[0], sourceDir).Run(); err != nil {
		fatal(fmt.Errorf("clone %s: %w", args[0], err))
	}
	if err := syncer.ImportSkills(sourceDir, filepath.Join(opts.BaseDir, "skills")); err != nil {
		fatal(err)
	}
	runSync(opts)
}

func runSync(opts syncer.Options) {
	plan, err := syncer.BuildPlan(opts)
	if err != nil {
		fatal(err)
	}
	if err := syncer.Apply(plan, opts, prompt); err != nil {
		fatal(err)
	}
	fmt.Println("done")
}

func prompt(action syncer.Action) (syncer.PromptDecision, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s %s -> %s. Overwrite? [y]es/[n]o/[a]ll/[q]uit: ", action.Type, action.Source, action.Dest)
		input, err := reader.ReadString('\n')
		if err != nil {
			return syncer.DecisionQuit, err
		}
		switch strings.ToLower(strings.TrimSpace(input)) {
		case "y", "yes":
			return syncer.DecisionYes, nil
		case "n", "no", "":
			return syncer.DecisionNo, nil
		case "a", "all":
			return syncer.DecisionAll, nil
		case "q", "quit":
			return syncer.DecisionQuit, nil
		default:
			fmt.Println("invalid choice")
		}
	}
}

func findMissingSources(opts syncer.Options) []string {
	var missing []string
	for _, target := range syncer.Targets(opts) {
		for _, path := range targetInputs(target) {
			if _, err := os.Stat(path); err != nil {
				if os.IsNotExist(err) {
					missing = append(missing, path)
				}
			}
		}
	}
	return missing
}

func targetInputs(target syncer.Target) []string {
	if len(target.Sources) > 0 {
		return target.Sources
	}
	if target.Source == "" {
		return nil
	}
	return []string{target.Source}
}

func defaultCodexHome() string {
	if value := os.Getenv("CODEX_HOME"); value != "" {
		return value
	}
	return "~/.codex"
}

func expandPath(path string) string {
	if path == "" {
		return path
	}
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, strings.TrimPrefix(path, "~"))
		}
	}
	return path
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
