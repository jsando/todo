package internal

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestPrintTable(t *testing.T) {
	issues := []*Issue{
		{ID: "aaaa1111", Type: "task", Priority: 2, Status: "open", Title: "First"},
		{ID: "bbbb2222", Type: "bug", Priority: 0, Status: "in_progress", Title: "Second", Labels: []string{"urgent"}},
	}
	var buf bytes.Buffer
	PrintTable(&buf, issues)
	out := buf.String()

	if !strings.Contains(out, "aaaa1111") {
		t.Error("missing first issue ID")
	}
	if !strings.Contains(out, "bbbb2222") {
		t.Error("missing second issue ID")
	}
	if !strings.Contains(out, "urgent") {
		t.Error("missing label")
	}
	// Header row
	if !strings.Contains(out, "ID") || !strings.Contains(out, "TYPE") {
		t.Error("missing header")
	}
}

func TestPrintJSON(t *testing.T) {
	issue := &Issue{ID: "aaaa1111", Title: "Test", Type: "task", Status: "open"}
	var buf bytes.Buffer
	PrintJSON(&buf, issue)

	var parsed map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if parsed["id"] != "aaaa1111" {
		t.Errorf("id = %v", parsed["id"])
	}
}

func TestPrintStatusSummary(t *testing.T) {
	issues := []*Issue{
		{Status: "open", Type: "task", Priority: 2},
		{Status: "open", Type: "bug", Priority: 0},
		{Status: "done", Type: "task", Priority: 2},
	}
	var buf bytes.Buffer
	PrintStatusSummary(&buf, issues)
	out := buf.String()

	if !strings.Contains(out, "Total: 3") {
		t.Errorf("missing total, got:\n%s", out)
	}
	if !strings.Contains(out, "open") || !strings.Contains(out, "done") {
		t.Error("missing status counts")
	}
}
