package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jsando/todo/internal"
)

func Search(args []string) error {
	fs := flag.NewFlagSet("search", flag.ExitOnError)
	jsonOut := fs.Bool("json", false, "JSON output")
	fs.Parse(reorderArgs(args))

	if fs.NArg() < 1 {
		return fmt.Errorf("usage: todo search <query>")
	}
	query := strings.ToLower(fs.Arg(0))

	todoDir, err := internal.FindRoot()
	if err != nil {
		return err
	}

	issues, err := internal.LoadEverything(todoDir)
	if err != nil {
		return err
	}

	var matches []*internal.Issue
	for _, i := range issues {
		if strings.Contains(strings.ToLower(i.Title), query) ||
			strings.Contains(strings.ToLower(i.Body), query) {
			matches = append(matches, i)
		}
	}

	if *jsonOut {
		internal.PrintJSON(os.Stdout, matches)
	} else {
		internal.PrintTable(os.Stdout, matches)
	}
	return nil
}
