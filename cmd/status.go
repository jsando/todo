package cmd

import (
	"flag"
	"os"

	"github.com/jsando/todo/internal"
)

func Status(args []string) error {
	fs := flag.NewFlagSet("status", flag.ExitOnError)
	jsonOut := fs.Bool("json", false, "JSON output")
	fs.Parse(args)

	todoDir, err := internal.FindRoot()
	if err != nil {
		return err
	}

	issues, err := internal.LoadEverything(todoDir)
	if err != nil {
		return err
	}

	if *jsonOut {
		summary := map[string]interface{}{
			"total":  len(issues),
			"issues": issues,
		}
		internal.PrintJSON(os.Stdout, summary)
	} else {
		internal.PrintStatusSummary(os.Stdout, issues)
	}
	return nil
}
