package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jsando/todo/internal"
)

func Update(args []string) error {
	fs := flag.NewFlagSet("update", flag.ExitOnError)
	status := fs.String("status", "", "New status")
	priority := fs.Int("priority", -1, "New priority")
	title := fs.String("title", "", "New title")
	addDep := fs.String("add-dep", "", "Add dependency")
	rmDep := fs.String("rm-dep", "", "Remove dependency")
	var addLabels stringList
	var rmLabels stringList
	fs.Var(&addLabels, "add-label", "Add label (repeatable)")
	fs.Var(&rmLabels, "rm-label", "Remove label (repeatable)")
	epic := fs.String("epic", "", "Set epic")
	jsonOut := fs.Bool("json", false, "JSON output")
	fs.Parse(reorderArgs(args))

	if fs.NArg() < 1 {
		return fmt.Errorf("usage: todo update <id> [flags]")
	}

	todoDir, err := internal.FindRoot()
	if err != nil {
		return err
	}

	issue, oldPath, err := internal.Resolve(todoDir, fs.Arg(0))
	if err != nil {
		return err
	}

	oldFilename := issue.Filename()

	if *status != "" {
		issue.Status = *status
	}
	if *priority >= 0 {
		issue.Priority = *priority
	}
	if *title != "" {
		issue.Title = *title
	}
	if *addDep != "" {
		issue.Deps = append(issue.Deps, *addDep)
	}
	if *rmDep != "" {
		issue.Deps = removeStr(issue.Deps, *rmDep)
	}
	for _, l := range addLabels {
		issue.Labels = append(issue.Labels, l)
	}
	for _, l := range rmLabels {
		issue.Labels = removeStr(issue.Labels, l)
	}
	if *epic != "" {
		issue.Epic = *epic
	}
	issue.Updated = internal.Today()

	// If filename changed (title changed), remove old file
	dir := filepath.Dir(oldPath)
	if issue.Filename() != oldFilename {
		os.Remove(oldPath)
	}

	if err := internal.WriteIssue(dir, issue); err != nil {
		return err
	}
	issue.Path = filepath.Join(dir, issue.Filename())

	if *jsonOut {
		internal.PrintJSON(os.Stdout, issue)
	} else {
		fmt.Printf("Updated %s\n", issue.ID)
	}
	return nil
}

func removeStr(slice []string, val string) []string {
	var result []string
	for _, s := range slice {
		if s != val {
			result = append(result, s)
		}
	}
	return result
}
