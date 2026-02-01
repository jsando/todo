package cmd

import "strings"

// reorderArgs moves flag arguments before positional arguments so that
// Go's flag package can parse them correctly. It recognizes flags as args
// starting with "-" and handles flags with separate values (e.g. -t feature).
// Boolean flags known to not take values: --json
var boolFlags = map[string]bool{
	"--json": true,
}

func reorderArgs(args []string) []string {
	var flags []string
	var positional []string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			positional = append(positional, args[i+1:]...)
			break
		}
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
			// If it's not a boolean flag and next arg exists and doesn't start with -, consume it
			if !boolFlags[arg] && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				i++
				flags = append(flags, args[i])
			}
		} else {
			positional = append(positional, arg)
		}
	}
	return append(flags, positional...)
}
