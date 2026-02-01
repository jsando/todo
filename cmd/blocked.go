package cmd

import (
	"flag"
	"os"

	"github.com/jsando/todo/internal"
)

func Blocked(args []string) error {
	fs := flag.NewFlagSet("blocked", flag.ExitOnError)
	jsonOut := fs.Bool("json", false, "JSON output")
	fs.Parse(args)

	todoDir, err := internal.FindRoot()
	if err != nil {
		return err
	}

	issues, err := internal.LoadAll(todoDir)
	if err != nil {
		return err
	}

	archived, _ := internal.LoadArchived(todoDir)
	doneIDs := map[string]bool{}
	for _, a := range archived {
		doneIDs[a.ID] = true
	}

	var blocked []*internal.Issue
	for _, i := range issues {
		if i.Status != "open" && i.Status != "in_progress" {
			continue
		}
		for _, dep := range i.Deps {
			if !doneIDs[dep] {
				blocked = append(blocked, i)
				break
			}
		}
	}

	if *jsonOut {
		internal.PrintJSON(os.Stdout, blocked)
	} else {
		internal.PrintTable(os.Stdout, blocked)
	}
	return nil
}
