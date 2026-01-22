package syncer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type PromptFunc func(action Action) (PromptDecision, error)

type PromptDecision int

const (
	DecisionYes PromptDecision = iota
	DecisionNo
	DecisionAll
	DecisionQuit
)

func Apply(plan Plan, opts Options, prompt PromptFunc) error {
	applyAll := false
	for _, action := range plan.Actions {
		if action.Type == ActionNoop {
			continue
		}
		if needsConfirmation(action) && !applyAll {
			decision, err := prompt(action)
			if err != nil {
				return err
			}
			switch decision {
			case DecisionQuit:
				return fmt.Errorf("aborted")
			case DecisionNo:
				continue
			case DecisionAll:
				applyAll = true
			case DecisionYes:
				// proceed
			}
		}

		if action.Type == ActionDelete {
			if err := os.Remove(action.Dest); err != nil && !os.IsNotExist(err) {
				return err
			}
			continue
		}

		if opts.UseSymlink {
			if err := replaceSymlink(action.Source, action.Dest); err != nil {
				return err
			}
			continue
		}

		if err := copyFile(action.Source, action.Dest); err != nil {
			return err
		}
	}
	return nil
}

func needsConfirmation(action Action) bool {
	if action.Type == ActionDelete && action.Note == "remove empty GEMINI.md" {
		return false
	}
	return action.Type == ActionUpdate || action.Type == ActionDelete
}

func replaceSymlink(src, dest string) error {
	if err := os.RemoveAll(dest); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	return os.Symlink(src, dest)
}

func copyFile(src, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return err
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	return os.Chmod(dest, info.Mode())
}
