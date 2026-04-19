package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jsando/todo/internal"
)

func Create(args []string) error {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	typ := fs.String("type", "task", "Issue type (task|bug|feature|epic|chore)")
	priority := fs.Int("priority", 2, "Priority (0-4)")
	epic := fs.String("epic", "", "Parent epic ID")
	dep := fs.String("dep", "", "Dependency issue ID")
	body := fs.String("body", "", "Markdown body (use '-' to read from stdin)")
	var labels stringList
	fs.Var(&labels, "label", "Label (repeatable)")
	jsonOut := fs.Bool("json", false, "JSON output")
	fs.Parse(reorderArgs(args))

	if fs.NArg() < 1 {
		return fmt.Errorf("usage: todo create \"title\" [--type X] [--priority X] [--epic id] [--dep id] [--label X] [--body X]")
	}
	title := fs.Arg(0)

	todoDir, err := internal.FindRoot()
	if err != nil {
		return err
	}

	bodyText := *body
	if bodyText == "-" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("reading body from stdin: %w", err)
		}
		bodyText = string(data)
	}

	issue := &internal.Issue{
		ID:       internal.NewID(),
		Title:    title,
		Type:     *typ,
		Status:   "open",
		Priority: *priority,
		Epic:     *epic,
		Created:  internal.Today(),
		Updated:  internal.Today(),
		Labels:   []string(labels),
		Body:     bodyText,
	}
	if *dep != "" {
		issue.Deps = []string{*dep}
	}

	if err := internal.WriteIssue(todoDir, issue); err != nil {
		return err
	}
	issue.Path = filepath.Join(todoDir, issue.Filename())

	if *jsonOut {
		internal.PrintJSON(os.Stdout, issue)
	} else {
		fmt.Println(issue.ID)
	}
	return nil
}
