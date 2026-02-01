package main

import (
	"fmt"
	"os"

	"github.com/jsando/todo/cmd"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
	command := os.Args[1]
	args := os.Args[2:]

	if command == "help" || command == "--help" || command == "-h" {
		if len(args) > 0 {
			printCommandHelp(args[0])
		} else {
			printUsage()
		}
		return
	}

	commands := map[string]func([]string) error{
		"init":    cmd.Init,
		"create":  cmd.Create,
		"list":    cmd.List,
		"show":    cmd.Show,
		"update":  cmd.Update,
		"close":   cmd.Close,
		"reopen":  cmd.Reopen,
		"ready":   cmd.Ready,
		"blocked": cmd.Blocked,
		"status":  cmd.Status,
		"search":  cmd.Search,
	}

	fn, ok := commands[command]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
	if err := fn(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Print(`todo — Markdown-based issue tracker

Usage: todo <command> [args]

Commands:
  init              Create .todo/ directory
  create "title"    Create a new issue
  list              List issues (with filters)
  show <id>         Show full issue details
  update <id>       Update issue fields
  close <id>        Close issue and move to archive
  reopen <id>       Reopen an archived issue
  ready             Show issues ready to work on (all deps resolved)
  blocked           Show issues waiting on dependencies
  status            Summary counts by status, type, and priority
  search <query>    Search issue titles and descriptions

All query commands accept --json for machine-readable output.
Partial IDs work like git: "todo show a3f7" matches "a3f7c201".

Run "todo help <command>" for details and examples.

Quick start:
  todo init
  todo create "Set up CI pipeline" --type task --priority 1
  todo create "Fix login timeout" --type bug --priority 0
  todo list
  todo ready
  todo update a3f7 --status in_progress
  todo close a3f7 --reason "Deployed to prod"
`)
}

var commandHelp = map[string]string{
	"init": `todo init

Create the .todo/ and .todo/archive/ directories in the current folder.

Example:
  todo init
`,

	"create": `todo create "title" [flags]

Create a new issue and print its ID.

Flags:
  --type      Issue type: task, bug, feature, epic, chore (default: task)
  --priority  Priority: 0 (critical) to 4 (backlog) (default: 2)
  --epic      Parent epic issue ID
  --dep       Dependency issue ID (this issue is blocked by it)
  --json      Print created issue as JSON

Examples:
  todo create "Add user authentication"
  todo create "Fix memory leak in worker" --type bug --priority 0
  todo create "Design new dashboard" --type feature --priority 1
  todo create "Refactor config loading" --type chore --json

  Create an epic with child issues:
    todo create "User auth system" --type epic --priority 1
    # Use the printed epic ID (e.g. a3f7c201) to link child issues:
    todo create "Design auth API" --type task --epic a3f7c201
    todo create "Implement JWT tokens" --type task --epic a3f7c201
    todo create "Write auth tests" --type task --epic a3f7c201

  Create issues with dependencies (blocked-by):
    todo create "Deploy to prod" --dep a3f7c201
    # This issue won't appear in "todo ready" until a3f7c201 is closed.
`,

	"list": `todo list [flags]

List issues in a table. Filters can be combined.

Flags:
  --status    Filter by status: open, in_progress, done, cancelled
  --type      Filter by type: task, bug, feature, epic, chore
  --priority  Filter by priority: 0-4
  --epic      Filter by parent epic ID
  --label     Filter by label
  --json      Output as JSON array

Examples:
  todo list
  todo list --status open
  todo list --type bug --priority 0
  todo list --status in_progress --json
  todo list --label auth

  List all issues belonging to an epic:
    todo list --epic a3f7c201

  List all epics:
    todo list --type epic
`,

	"show": `todo show <id> [--json]

Print the full issue file contents. Partial IDs work.

Examples:
  todo show a3f7c201
  todo show a3f7
  todo show a3f7 --json
`,

	"update": `todo update <id> [flags]

Update one or more fields on an issue. Only specified fields are changed.

Flags:
  --status      Set status: open, in_progress, done, cancelled
  --priority    Set priority: 0-4
  --title       Set new title
  --epic        Set parent epic ID
  --add-dep     Add a dependency
  --rm-dep      Remove a dependency
  --add-label   Add a label
  --rm-label    Remove a label
  --json        Print updated issue as JSON

Examples:
  todo update a3f7 --status in_progress
  todo update a3f7 --priority 0 --add-label urgent
  todo update a3f7 --add-dep b912de44
  todo update a3f7 --rm-label wontfix
  todo update a3f7 --title "Revised title" --json
`,

	"close": `todo close <id> [--reason X] [--json]

Set status to done and move the issue file to .todo/archive/.

Flags:
  --reason    Append a closing reason to the issue body
  --json      Print closed issue as JSON

Examples:
  todo close a3f7
  todo close a3f7 --reason "Deployed to production"
  todo close a3f7 --json
`,

	"reopen": `todo reopen <id> [--json]

Move an issue from .todo/archive/ back to .todo/ and set status to open.

Examples:
  todo reopen a3f7
  todo reopen a3f7 --json
`,

	"ready": `todo ready [--json]

List open/in_progress issues where all dependencies are resolved (closed).
Issues with no dependencies are always included.

Examples:
  todo ready
  todo ready --json
`,

	"blocked": `todo blocked [--json]

List open/in_progress issues that have at least one unresolved dependency.

Examples:
  todo blocked
  todo blocked --json
`,

	"status": `todo status [--json]

Show summary counts grouped by status, type, and priority.
Includes both active and archived issues.

Examples:
  todo status
  todo status --json
`,

	"search": `todo search <query> [--json]

Case-insensitive search across issue titles and descriptions.
Searches both active and archived issues.

Examples:
  todo search auth
  todo search "login bug"
  todo search timeout --json
`,
}

func printCommandHelp(name string) {
	help, ok := commandHelp[name]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", name)
		os.Exit(1)
	}
	fmt.Print(help)
}
