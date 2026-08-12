package syncer

import (
	"fmt"
	"os"
	"path/filepath"
)

type ActionType string

const (
	ActionCreate ActionType = "create"
	ActionUpdate ActionType = "update"
	ActionDelete ActionType = "delete"
	ActionNoop   ActionType = "noop"
)

type Action struct {
	Type    ActionType
	Source  string
	Dest    string
	Note    string
	Content []byte
}

type Plan struct {
	Actions []Action
}

func BuildPlan(opts Options) (Plan, error) {
	var actions []Action
	if action, ok := emptyGeminiCleanup(opts.GeminiHome); ok {
		actions = append(actions, action)
	}

	for _, target := range Targets(opts) {
		if !isManagedTarget(target) {
			linkAction, err := buildLinkAction(target)
			if err != nil {
				return Plan{}, err
			}
			actions = append(actions, linkAction)
			continue
		}

		switch target.Kind {
		case KindFile, KindDir:
			return Plan{}, fmt.Errorf("unmanaged target %s must be linked", target.Name)
		case KindGeminiConfig, KindCodexConfig, KindAgyMcpConfig, KindOpenCodeConfig:
			action, err := buildManagedAction(target, opts)
			if err != nil {
				return Plan{}, err
			}
			actions = append(actions, action)
		default:
			return Plan{}, fmt.Errorf("unknown target kind: %s", target.Kind)
		}
	}

	return Plan{Actions: actions}, nil
}

func isManagedTarget(target Target) bool {
	return target.Kind == KindGeminiConfig || target.Kind == KindCodexConfig || target.Kind == KindAgyMcpConfig || target.Kind == KindOpenCodeConfig
}

func emptyGeminiCleanup(geminiHome string) (Action, bool) {
	path := filepath.Join(geminiHome, "GEMINI.md")
	info, err := os.Stat(path)
	if err != nil {
		return Action{}, false
	}
	if info.IsDir() {
		return Action{}, false
	}
	if info.Size() != 0 {
		return Action{}, false
	}
	return Action{
		Type: ActionDelete,
		Dest: path,
		Note: "remove empty GEMINI.md",
	}, true
}

func buildLinkAction(target Target) (Action, error) {
	info, err := os.Lstat(target.Dest)
	if err != nil {
		if os.IsNotExist(err) {
			return Action{Type: ActionCreate, Source: target.Source, Dest: target.Dest, Note: "symlink"}, nil
		}
		return Action{}, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		link, err := os.Readlink(target.Dest)
		if err != nil {
			return Action{}, err
		}
		if filepath.Clean(link) == filepath.Clean(target.Source) {
			return Action{Type: ActionNoop, Source: target.Source, Dest: target.Dest, Note: "symlink"}, nil
		}
	}
	return Action{Type: ActionUpdate, Source: target.Source, Dest: target.Dest, Note: "symlink"}, nil
}

func isSymlink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink != 0
}
