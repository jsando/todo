package cmd

import (
	"flag"
	"os"

	"github.com/jsando/todo/internal"
)

func Ready(args []string) error {
	fs := flag.NewFlagSet("ready", flag.ExitOnError)
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

	// Build set of done issue IDs (from archive)
	archived, _ := internal.LoadArchived(todoDir)
	doneIDs := map[string]bool{}
	for _, a := range archived {
		doneIDs[a.ID] = true
	}

	var ready []*internal.Issue
	for _, i := range issues {
		if i.Status != "open" && i.Status != "in_progress" {
			continue
		}
		allResolved := true
		for _, dep := range i.Deps {
			if !doneIDs[dep] {
				allResolved = false
				break
			}
		}
		if allResolved {
			ready = append(ready, i)
		}
	}

	if *jsonOut {
		internal.PrintJSON(os.Stdout, ready)
	} else {
		internal.PrintTable(os.Stdout, ready)
	}
	return nil
}
