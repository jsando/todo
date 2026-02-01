package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/jsando/todo/internal"
)

func Show(args []string) error {
	fs := flag.NewFlagSet("show", flag.ExitOnError)
	jsonOut := fs.Bool("json", false, "JSON output")
	fs.Parse(reorderArgs(args))

	if fs.NArg() < 1 {
		return fmt.Errorf("usage: todo show <id>")
	}

	todoDir, err := internal.FindRoot()
	if err != nil {
		return err
	}

	issue, path, err := internal.Resolve(todoDir, fs.Arg(0))
	if err != nil {
		return err
	}

	if *jsonOut {
		internal.PrintJSON(os.Stdout, issue)
	} else {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fmt.Print(string(data))
	}
	return nil
}
