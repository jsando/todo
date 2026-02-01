package internal

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Issue struct {
	ID       string   `yaml:"id" json:"id"`
	Title    string   `yaml:"title" json:"title"`
	Type     string   `yaml:"type" json:"type"`
	Status   string   `yaml:"status" json:"status"`
	Priority int      `yaml:"priority" json:"priority"`
	Epic     string   `yaml:"epic,omitempty" json:"epic,omitempty"`
	Deps     []string `yaml:"deps,omitempty" json:"deps,omitempty"`
	Created  string   `yaml:"created" json:"created"`
	Updated  string   `yaml:"updated" json:"updated"`
	Labels   []string `yaml:"labels,omitempty" json:"labels,omitempty"`

	Body string `yaml:"-" json:"body,omitempty"`
}

func NewID() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%08x", b)
}

func Today() string {
	return time.Now().Format("2006-01-02")
}

func Slugify(title string) string {
	s := strings.ToLower(title)
	s = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			return r
		}
		return '-'
	}, s)
	// collapse multiple dashes
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	s = strings.Trim(s, "-")
	if len(s) > 50 {
		s = s[:50]
		s = strings.TrimRight(s, "-")
	}
	return s
}

func (i *Issue) Filename() string {
	return i.ID + "-" + Slugify(i.Title) + ".md"
}

func ParseIssue(data []byte) (*Issue, error) {
	content := string(data)
	if !strings.HasPrefix(content, "---\n") {
		return nil, fmt.Errorf("missing YAML front matter")
	}
	end := strings.Index(content[4:], "\n---\n")
	if end < 0 {
		return nil, fmt.Errorf("missing end of YAML front matter")
	}
	yamlPart := content[4 : 4+end]
	body := content[4+end+5:]

	var issue Issue
	if err := yaml.Unmarshal([]byte(yamlPart), &issue); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}
	issue.Body = strings.TrimSpace(body)
	return &issue, nil
}

func (i *Issue) Serialize() []byte {
	yamlBytes, _ := yaml.Marshal(i)
	var sb strings.Builder
	sb.WriteString("---\n")
	sb.Write(yamlBytes)
	sb.WriteString("---\n")
	if i.Body != "" {
		sb.WriteString("\n")
		sb.WriteString(i.Body)
		sb.WriteString("\n")
	}
	return []byte(sb.String())
}
