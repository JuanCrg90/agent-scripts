package syncer

import (
	"fmt"
	"io/fs"
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
		if opts.UseSymlink && !isManagedTarget(target) {
			linkAction, err := buildLinkAction(target)
			if err != nil {
				return Plan{}, err
			}
			actions = append(actions, linkAction)
			continue
		}

		switch target.Kind {
		case KindFile:
			actions = append(actions, buildFileAction(target.Source, target.Dest))
		case KindDir:
			dirActions, err := buildDirActions(target.Source, target.Dest)
			if err != nil {
				return Plan{}, err
			}
			actions = append(actions, dirActions...)
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

func buildFileAction(src, dest string) Action {
	if isSymlink(dest) {
		return Action{Type: ActionUpdate, Source: src, Dest: dest, Note: "replace symlink"}
	}
	if sameFile(src, dest) {
		return Action{Type: ActionNoop, Source: src, Dest: dest}
	}
	if _, err := os.Stat(dest); err != nil {
		if os.IsNotExist(err) {
			return Action{Type: ActionCreate, Source: src, Dest: dest}
		}
		return Action{Type: ActionUpdate, Source: src, Dest: dest, Note: err.Error()}
	}
	return Action{Type: ActionUpdate, Source: src, Dest: dest}
}

func buildDirActions(srcDir, destDir string) ([]Action, error) {
	var actions []Action
	err := filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		dest := filepath.Join(destDir, rel)
		actions = append(actions, buildFileAction(path, dest))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return actions, nil
}

func sameFile(src, dest string) bool {
	srcInfo, err := os.Stat(src)
	if err != nil || srcInfo.IsDir() {
		return false
	}
	destInfo, err := os.Stat(dest)
	if err != nil || destInfo.IsDir() {
		return false
	}
	if srcInfo.Size() != destInfo.Size() {
		return false
	}
	srcBytes, err := os.ReadFile(src)
	if err != nil {
		return false
	}
	destBytes, err := os.ReadFile(dest)
	if err != nil {
		return false
	}
	return string(srcBytes) == string(destBytes)
}

func isSymlink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink != 0
}
