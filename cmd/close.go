package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jsando/todo/internal"
)

func Close(args []string) error {
	fs := flag.NewFlagSet("close", flag.ExitOnError)
	reason := fs.String("reason", "", "Reason for closing")
	jsonOut := fs.Bool("json", false, "JSON output")
	fs.Parse(reorderArgs(args))

	if fs.NArg() < 1 {
		return fmt.Errorf("usage: todo close <id> [--reason X]")
	}

	todoDir, err := internal.FindRoot()
	if err != nil {
		return err
	}

	issue, oldPath, err := internal.Resolve(todoDir, fs.Arg(0))
	if err != nil {
		return err
	}

	issue.Status = "done"
	issue.Updated = internal.Today()
	if *reason != "" {
		if issue.Body != "" {
			issue.Body += "\n\n"
		}
		issue.Body += "Closed: " + *reason
	}

	// Move to archive
	archiveDir := filepath.Join(todoDir, "archive")
	if err := internal.WriteIssue(archiveDir, issue); err != nil {
		return err
	}
	os.Remove(oldPath)

	if *jsonOut {
		internal.PrintJSON(os.Stdout, issue)
	} else {
		fmt.Printf("Closed %s\n", issue.ID)
	}
	return nil
}
