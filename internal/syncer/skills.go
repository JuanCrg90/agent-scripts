package syncer

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

// ImportSkills copies skills from a checked-out repository into the canonical
// skill store. Existing skills must be identical; conflicts are never merged
// or overwritten automatically.
func ImportSkills(repoDir, skillsDir string) error {
	sourceDir := filepath.Join(repoDir, "skills")
	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return fmt.Errorf("read imported skills: %w", err)
	}

	var imports []string
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == ".system" {
			continue
		}

		source := filepath.Join(sourceDir, entry.Name())
		dest := filepath.Join(skillsDir, entry.Name())
		equal, exists, err := sameDir(source, dest)
		if err != nil {
			return err
		}
		if exists && !equal {
			return fmt.Errorf("skill conflict: %s already exists at %s", entry.Name(), dest)
		}
		if !exists {
			imports = append(imports, entry.Name())
		}
	}

	for _, name := range imports {
		source := filepath.Join(sourceDir, name)
		dest := filepath.Join(skillsDir, name)
		if err := copyDir(source, dest); err != nil {
			return fmt.Errorf("import skill %s: %w", name, err)
		}
	}
	return nil
}

func sameDir(source, dest string) (equal, exists bool, err error) {
	if _, err := os.Stat(dest); err != nil {
		if os.IsNotExist(err) {
			return false, false, nil
		}
		return false, false, err
	}

	var sourceFiles, destFiles []string
	if err := filepath.WalkDir(source, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			rel, err := filepath.Rel(source, path)
			if err != nil {
				return err
			}
			sourceFiles = append(sourceFiles, rel)
		}
		return nil
	}); err != nil {
		return false, true, err
	}
	if err := filepath.WalkDir(dest, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			rel, err := filepath.Rel(dest, path)
			if err != nil {
				return err
			}
			destFiles = append(destFiles, rel)
		}
		return nil
	}); err != nil {
		return false, true, err
	}
	sort.Strings(sourceFiles)
	sort.Strings(destFiles)
	if len(sourceFiles) != len(destFiles) {
		return false, true, nil
	}
	for index, rel := range sourceFiles {
		if rel != destFiles[index] || !sameFile(filepath.Join(source, rel), filepath.Join(dest, rel)) {
			return false, true, nil
		}
	}
	return true, true, nil
}

func sameFile(source, dest string) bool {
	sourceBytes, err := os.ReadFile(source)
	if err != nil {
		return false
	}
	destBytes, err := os.ReadFile(dest)
	if err != nil {
		return false
	}
	return bytes.Equal(sourceBytes, destBytes)
}

func copyDir(source, dest string) error {
	return filepath.WalkDir(source, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dest, rel)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target)
	})
}
