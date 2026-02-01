package internal

import (
	"strings"
	"testing"
)

func TestNewID(t *testing.T) {
	id := NewID()
	if len(id) != 8 {
		t.Fatalf("expected 8-char ID, got %q (len %d)", id, len(id))
	}
	// Should be hex
	for _, c := range id {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Fatalf("non-hex character %c in ID %q", c, id)
		}
	}
	// Two IDs should differ
	id2 := NewID()
	if id == id2 {
		t.Fatalf("two generated IDs are identical: %s", id)
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"Add JWT authentication", "add-jwt-authentication"},
		{"Fix login bug!", "fix-login-bug"},
		{"  spaces  everywhere  ", "spaces-everywhere"},
		{"UPPERCASE", "uppercase"},
		{"a/b/c", "a-b-c"},
		{"", ""},
		{strings.Repeat("a", 100), strings.Repeat("a", 50)},
	}
	for _, tt := range tests {
		got := Slugify(tt.input)
		if got != tt.want {
			t.Errorf("Slugify(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestParseAndSerialize(t *testing.T) {
	input := `---
id: a3f7c201
title: Add JWT authentication
type: feature
status: open
priority: 1
epic: c44b2e10
deps:
    - b912de44
created: "2026-02-01"
updated: "2026-02-01"
labels:
    - auth
    - api
---

This is the description.
`
	issue, err := ParseIssue([]byte(input))
	if err != nil {
		t.Fatalf("ParseIssue: %v", err)
	}
	if issue.ID != "a3f7c201" {
		t.Errorf("ID = %q, want %q", issue.ID, "a3f7c201")
	}
	if issue.Title != "Add JWT authentication" {
		t.Errorf("Title = %q", issue.Title)
	}
	if issue.Type != "feature" {
		t.Errorf("Type = %q", issue.Type)
	}
	if issue.Status != "open" {
		t.Errorf("Status = %q", issue.Status)
	}
	if issue.Priority != 1 {
		t.Errorf("Priority = %d", issue.Priority)
	}
	if issue.Epic != "c44b2e10" {
		t.Errorf("Epic = %q", issue.Epic)
	}
	if len(issue.Deps) != 1 || issue.Deps[0] != "b912de44" {
		t.Errorf("Deps = %v", issue.Deps)
	}
	if len(issue.Labels) != 2 || issue.Labels[0] != "auth" || issue.Labels[1] != "api" {
		t.Errorf("Labels = %v", issue.Labels)
	}
	if issue.Body != "This is the description." {
		t.Errorf("Body = %q", issue.Body)
	}

	// Round-trip: serialize and re-parse
	data := issue.Serialize()
	issue2, err := ParseIssue(data)
	if err != nil {
		t.Fatalf("re-parse: %v", err)
	}
	if issue2.ID != issue.ID || issue2.Title != issue.Title || issue2.Body != issue.Body {
		t.Errorf("round-trip mismatch: got ID=%q Title=%q Body=%q", issue2.ID, issue2.Title, issue2.Body)
	}
}

func TestParseIssue_NoBody(t *testing.T) {
	input := `---
id: abcd1234
title: No body issue
type: task
status: open
priority: 2
created: "2026-02-01"
updated: "2026-02-01"
---
`
	issue, err := ParseIssue([]byte(input))
	if err != nil {
		t.Fatalf("ParseIssue: %v", err)
	}
	if issue.Body != "" {
		t.Errorf("Body = %q, want empty", issue.Body)
	}
}

func TestParseIssue_Errors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"no front matter", "just some text"},
		{"no closing delimiter", "---\nid: abc\n"},
	}
	for _, tt := range tests {
		_, err := ParseIssue([]byte(tt.input))
		if err == nil {
			t.Errorf("%s: expected error, got nil", tt.name)
		}
	}
}

func TestValidate_MissingFields(t *testing.T) {
	input := `---
id: abcd1234
title: ""
type: task
status: open
priority: 2
created: "2026-02-01"
updated: "2026-02-01"
---
`
	_, err := ParseIssue([]byte(input))
	if err == nil {
		t.Fatal("expected validation error for empty title")
	}
	if !strings.Contains(err.Error(), "missing required field: title") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_InvalidType(t *testing.T) {
	input := `---
id: abcd1234
title: Test
type: banana
status: open
priority: 2
created: "2026-02-01"
updated: "2026-02-01"
---
`
	_, err := ParseIssue([]byte(input))
	if err == nil {
		t.Fatal("expected validation error for invalid type")
	}
	if !strings.Contains(err.Error(), "invalid type") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_InvalidStatus(t *testing.T) {
	input := `---
id: abcd1234
title: Test
type: task
status: pending
priority: 2
created: "2026-02-01"
updated: "2026-02-01"
---
`
	_, err := ParseIssue([]byte(input))
	if err == nil {
		t.Fatal("expected validation error for invalid status")
	}
	if !strings.Contains(err.Error(), "invalid status") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_InvalidPriority(t *testing.T) {
	input := `---
id: abcd1234
title: Test
type: task
status: open
priority: 9
created: "2026-02-01"
updated: "2026-02-01"
---
`
	_, err := ParseIssue([]byte(input))
	if err == nil {
		t.Fatal("expected validation error for invalid priority")
	}
	if !strings.Contains(err.Error(), "invalid priority") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	input := `---
id: ""
title: ""
type: banana
status: nope
priority: 9
created: ""
updated: ""
---
`
	_, err := ParseIssue([]byte(input))
	if err == nil {
		t.Fatal("expected validation errors")
	}
	errStr := err.Error()
	for _, want := range []string{"id", "title", "type", "status", "priority", "created", "updated"} {
		if !strings.Contains(errStr, want) {
			t.Errorf("error missing mention of %q: %v", want, err)
		}
	}
}

func TestFilename(t *testing.T) {
	issue := &Issue{ID: "a3f7c201", Title: "Add JWT auth"}
	want := "a3f7c201-add-jwt-auth.md"
	if got := issue.Filename(); got != want {
		t.Errorf("Filename() = %q, want %q", got, want)
	}
}
