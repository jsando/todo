package cmd

import (
	"flag"
	"os"

	"github.com/jsando/todo/internal"
)

func List(args []string) error {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	status := fs.String("status", "", "Filter by status")
	typ := fs.String("type", "", "Filter by type")
	priority := fs.Int("priority", -1, "Filter by priority")
	epic := fs.String("epic", "", "Filter by epic")
	label := fs.String("label", "", "Filter by label")
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

	var filtered []*internal.Issue
	for _, i := range issues {
		if *status != "" && i.Status != *status {
			continue
		}
		if *typ != "" && i.Type != *typ {
			continue
		}
		if *priority >= 0 && i.Priority != *priority {
			continue
		}
		if *epic != "" && i.Epic != *epic {
			continue
		}
		if *label != "" && !hasLabel(i, *label) {
			continue
		}
		filtered = append(filtered, i)
	}

	if *jsonOut {
		internal.PrintJSON(os.Stdout, filtered)
	} else {
		internal.PrintTable(os.Stdout, filtered)
	}
	return nil
}

func hasLabel(i *internal.Issue, label string) bool {
	for _, l := range i.Labels {
		if l == label {
			return true
		}
	}
	return false
}
