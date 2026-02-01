package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jsando/todo/internal"
)

func Reopen(args []string) error {
	fs := flag.NewFlagSet("reopen", flag.ExitOnError)
	jsonOut := fs.Bool("json", false, "JSON output")
	fs.Parse(reorderArgs(args))

	if fs.NArg() < 1 {
		return fmt.Errorf("usage: todo reopen <id>")
	}

	todoDir, err := internal.FindRoot()
	if err != nil {
		return err
	}

	// Look in archive
	archiveDir := filepath.Join(todoDir, "archive")
	issue, oldPath, err := internal.Resolve(archiveDir, fs.Arg(0))
	if err != nil {
		return fmt.Errorf("issue not found in archive: %w", err)
	}

	issue.Status = "open"
	issue.Updated = internal.Today()

	if err := internal.WriteIssue(todoDir, issue); err != nil {
		return err
	}
	os.Remove(oldPath)

	if *jsonOut {
		internal.PrintJSON(os.Stdout, issue)
	} else {
		fmt.Printf("Reopened %s\n", issue.ID)
	}
	return nil
}
