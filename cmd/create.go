package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/jsando/todo/internal"
)

func Create(args []string) error {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	typ := fs.String("type", "task", "Issue type (task|bug|feature|epic|chore)")
	priority := fs.Int("priority", 2, "Priority (0-4)")
	epic := fs.String("epic", "", "Parent epic ID")
	dep := fs.String("dep", "", "Dependency issue ID")
	jsonOut := fs.Bool("json", false, "JSON output")
	fs.Parse(reorderArgs(args))

	if fs.NArg() < 1 {
		return fmt.Errorf("usage: todo create \"title\" [--type X] [--priority X] [--epic id] [--dep id]")
	}
	title := fs.Arg(0)

	todoDir, err := internal.FindRoot()
	if err != nil {
		return err
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
	}
	if *dep != "" {
		issue.Deps = []string{*dep}
	}

	if err := internal.WriteIssue(todoDir, issue); err != nil {
		return err
	}

	if *jsonOut {
		internal.PrintJSON(os.Stdout, issue)
	} else {
		fmt.Println(issue.ID)
	}
	return nil
}
