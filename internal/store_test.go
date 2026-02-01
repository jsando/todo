package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTodoDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	todoDir := filepath.Join(dir, ".todo")
	archiveDir := filepath.Join(todoDir, "archive")
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		t.Fatal(err)
	}
	return todoDir
}

func writeTestIssue(t *testing.T, dir string, issue *Issue) {
	t.Helper()
	if err := WriteIssue(dir, issue); err != nil {
		t.Fatal(err)
	}
}

func TestWriteAndLoadAll(t *testing.T) {
	todoDir := setupTodoDir(t)

	i1 := &Issue{ID: "aaaa1111", Title: "First", Type: "task", Status: "open", Priority: 2, Created: "2026-02-01", Updated: "2026-02-01"}
	i2 := &Issue{ID: "bbbb2222", Title: "Second", Type: "bug", Status: "open", Priority: 0, Created: "2026-02-01", Updated: "2026-02-01"}
	writeTestIssue(t, todoDir, i1)
	writeTestIssue(t, todoDir, i2)

	issues, err := LoadAll(todoDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 2 {
		t.Fatalf("LoadAll returned %d issues, want 2", len(issues))
	}
}

func TestResolve_FullID(t *testing.T) {
	todoDir := setupTodoDir(t)
	i := &Issue{ID: "aaaa1111", Title: "Test", Type: "task", Status: "open", Priority: 2, Created: "2026-02-01", Updated: "2026-02-01"}
	writeTestIssue(t, todoDir, i)

	found, _, err := Resolve(todoDir, "aaaa1111")
	if err != nil {
		t.Fatal(err)
	}
	if found.ID != "aaaa1111" {
		t.Errorf("Resolve got ID %q", found.ID)
	}
}

func TestResolve_PartialID(t *testing.T) {
	todoDir := setupTodoDir(t)
	i := &Issue{ID: "abcd5678", Title: "Partial test", Type: "task", Status: "open", Priority: 2, Created: "2026-02-01", Updated: "2026-02-01"}
	writeTestIssue(t, todoDir, i)

	found, _, err := Resolve(todoDir, "abcd")
	if err != nil {
		t.Fatal(err)
	}
	if found.ID != "abcd5678" {
		t.Errorf("Resolve got ID %q", found.ID)
	}
}

func TestResolve_Ambiguous(t *testing.T) {
	todoDir := setupTodoDir(t)
	writeTestIssue(t, todoDir, &Issue{ID: "aa001111", Title: "One", Type: "task", Status: "open", Priority: 2, Created: "2026-02-01", Updated: "2026-02-01"})
	writeTestIssue(t, todoDir, &Issue{ID: "aa002222", Title: "Two", Type: "task", Status: "open", Priority: 2, Created: "2026-02-01", Updated: "2026-02-01"})

	_, _, err := Resolve(todoDir, "aa00")
	if err == nil {
		t.Fatal("expected ambiguous error")
	}
}

func TestResolve_NotFound(t *testing.T) {
	todoDir := setupTodoDir(t)
	_, _, err := Resolve(todoDir, "zzzz")
	if err == nil {
		t.Fatal("expected not found error")
	}
}

func TestResolve_Archive(t *testing.T) {
	todoDir := setupTodoDir(t)
	archiveDir := filepath.Join(todoDir, "archive")
	i := &Issue{ID: "dead0000", Title: "Archived", Type: "task", Status: "done", Priority: 2, Created: "2026-02-01", Updated: "2026-02-01"}
	writeTestIssue(t, archiveDir, i)

	found, path, err := Resolve(todoDir, "dead")
	if err != nil {
		t.Fatal(err)
	}
	if found.ID != "dead0000" {
		t.Errorf("got ID %q", found.ID)
	}
	if filepath.Dir(path) != archiveDir {
		t.Errorf("expected archive path, got %q", path)
	}
}

func TestLoadEverything(t *testing.T) {
	todoDir := setupTodoDir(t)
	archiveDir := filepath.Join(todoDir, "archive")

	writeTestIssue(t, todoDir, &Issue{ID: "aaaa1111", Title: "Active", Type: "task", Status: "open", Priority: 2, Created: "2026-02-01", Updated: "2026-02-01"})
	writeTestIssue(t, archiveDir, &Issue{ID: "bbbb2222", Title: "Done", Type: "task", Status: "done", Priority: 2, Created: "2026-02-01", Updated: "2026-02-01"})

	all, err := LoadEverything(todoDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 {
		t.Fatalf("LoadEverything returned %d, want 2", len(all))
	}
}
