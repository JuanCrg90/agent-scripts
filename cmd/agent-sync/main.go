package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
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
	codexHome := fs.String("codex-home", defaultCodexHome(), "Codex home directory")
	geminiHome := fs.String("gemini-home", "~/.gemini", "Gemini home directory")
	useSymlink := fs.Bool("symlink", false, "Use symlinks instead of copying (single source of truth, but links can break if moved)")
	fs.Parse(os.Args[2:])

	opts := syncer.Options{
		BaseDir:    expandPath(*base),
		CodexHome:  expandPath(*codexHome),
		GeminiHome: expandPath(*geminiHome),
		UseSymlink: *useSymlink,
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
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Println("agent-sync <command> [flags]")
	fmt.Println("")
	fmt.Println("Commands: status, diff, plan, doctor, init, sync")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  --base <path>        Base directory (default: ~/Projects/agent-scripts)")
	fmt.Println("  --codex-home <path>  Codex home (default: $CODEX_HOME or ~/.codex)")
	fmt.Println("  --gemini-home <path> Gemini home (default: ~/.gemini)")
	fmt.Println("  --symlink            Use symlinks instead of copying. Pros: always in sync. Cons: links can break if you move the repo; some tools dislike symlinks.")
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
	if err := os.MkdirAll(filepath.Join(opts.GeminiHome, "antigravity"), 0o755); err != nil {
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
