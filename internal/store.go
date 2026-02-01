package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FindRoot walks up from cwd looking for .todo/ directory.
// Returns the path to the .todo/ directory.
func FindRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		todoDir := filepath.Join(dir, ".todo")
		if info, err := os.Stat(todoDir); err == nil && info.IsDir() {
			return todoDir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf(".todo directory not found (run 'todo init' first)")
		}
		dir = parent
	}
}

// LoadAll loads all issues from .todo/ (not archive).
func LoadAll(todoDir string) ([]*Issue, error) {
	return loadFrom(todoDir)
}

// LoadArchived loads all issues from .todo/archive/.
func LoadArchived(todoDir string) ([]*Issue, error) {
	archiveDir := filepath.Join(todoDir, "archive")
	return loadFrom(archiveDir)
}

// LoadEverything loads both active and archived issues.
func LoadEverything(todoDir string) ([]*Issue, error) {
	active, err := LoadAll(todoDir)
	if err != nil {
		return nil, err
	}
	archived, err := LoadArchived(todoDir)
	if err != nil {
		return nil, err
	}
	return append(active, archived...), nil
}

func loadFrom(dir string) ([]*Issue, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var issues []*Issue
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		issue, err := ParseIssue(data)
		if err != nil {
			continue // skip malformed files
		}
		issues = append(issues, issue)
	}
	return issues, nil
}

// Resolve finds an issue by full or partial ID. Returns the issue and its file path.
func Resolve(todoDir, partialID string) (*Issue, string, error) {
	// Search active issues
	issue, path, err := resolveIn(todoDir, partialID)
	if err == nil {
		return issue, path, nil
	}
	// Search archive
	archiveDir := filepath.Join(todoDir, "archive")
	issue, path, err = resolveIn(archiveDir, partialID)
	if err == nil {
		return issue, path, nil
	}
	return nil, "", fmt.Errorf("no issue matching '%s'", partialID)
}

func resolveIn(dir, partialID string) (*Issue, string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, "", err
	}
	var matches []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		if strings.HasPrefix(e.Name(), partialID) {
			matches = append(matches, e.Name())
		}
	}
	if len(matches) == 0 {
		return nil, "", fmt.Errorf("no match")
	}
	if len(matches) > 1 {
		return nil, "", fmt.Errorf("ambiguous ID '%s' matches %d issues", partialID, len(matches))
	}
	path := filepath.Join(dir, matches[0])
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}
	issue, err := ParseIssue(data)
	if err != nil {
		return nil, "", err
	}
	return issue, path, nil
}

// WriteIssue writes an issue to the given directory.
func WriteIssue(dir string, issue *Issue) error {
	path := filepath.Join(dir, issue.Filename())
	return os.WriteFile(path, issue.Serialize(), 0644)
}
