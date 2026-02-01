package cmd

import (
	"strings"
	"testing"
)

func TestReorderArgs(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want string
	}{
		{
			"flags before positional",
			[]string{"--type", "bug", "--priority", "0", "title"},
			"--type bug --priority 0 title",
		},
		{
			"title first then flags",
			[]string{"title", "--type", "bug", "--priority", "0"},
			"--type bug --priority 0 title",
		},
		{
			"mixed order",
			[]string{"--type", "bug", "title", "--priority", "0"},
			"--type bug --priority 0 title",
		},
		{
			"boolean flag",
			[]string{"title", "--json"},
			"--json title",
		},
		{
			"no flags",
			[]string{"just", "positional", "args"},
			"just positional args",
		},
		{
			"only flags",
			[]string{"--status", "open", "--json"},
			"--status open --json",
		},
	}
	for _, tt := range tests {
		got := strings.Join(reorderArgs(tt.in), " ")
		if got != tt.want {
			t.Errorf("%s: reorderArgs(%v) = %q, want %q", tt.name, tt.in, got, tt.want)
		}
	}
}
