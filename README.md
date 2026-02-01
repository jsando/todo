# todo

A lightweight CLI for project-level issue tracking using markdown files. Each issue is a separate `.md` file with YAML front matter in a `.todo/` directory. No database, no central index — just files.

## Install

```bash
go install github.com/jsando/todo@latest
```

Or build from source:

```bash
git clone https://github.com/jsando/todo.git
cd todo
go build -o todo .
```

## Quick Start

```bash
todo init
todo create "Set up CI pipeline" --type task --priority 1
todo create "Fix login timeout" --type bug --priority 0
todo list
todo update a3f7 --status in_progress
todo close a3f7 --reason "Deployed to prod"
```

## How It Works

Issues are stored as markdown files in `.todo/`:

```
.todo/
  a3f7c201-set-up-ci-pipeline.md
  b912de44-fix-login-timeout.md
  archive/
    c44b2e10-old-issue.md
```

Each file has YAML front matter:

```markdown
---
id: a3f7c201
title: Set up CI pipeline
type: task
status: open
priority: 1
created: 2026-02-01
updated: 2026-02-01
labels: [ci, infra]
---

Optional description in free-form markdown.
```

Closing an issue moves it to `.todo/archive/`. Reopening moves it back.

## Commands

| Command | Description |
|---|---|
| `todo init` | Create `.todo/` directory |
| `todo create "title"` | Create a new issue |
| `todo list` | List issues with optional filters |
| `todo show <id>` | Show full issue details |
| `todo update <id>` | Update issue fields |
| `todo close <id>` | Close and archive an issue |
| `todo reopen <id>` | Restore an archived issue |
| `todo ready` | Issues with all deps resolved |
| `todo blocked` | Issues waiting on dependencies |
| `todo status` | Summary counts |
| `todo search <query>` | Search titles and descriptions |

Run `todo help <command>` for flags and examples.

### Partial IDs

Like git, you don't need the full ID. `todo show a3f7` matches `a3f7c201-set-up-ci-pipeline.md`.

### JSON Output

All query commands accept `--json` for machine-readable output.

## Epics and Dependencies

Create an epic and link child issues:

```bash
todo create "User auth system" --type epic --priority 1
# prints: a3f7c201

todo create "Design auth API" --type task --epic a3f7c201
todo create "Implement JWT tokens" --type task --epic a3f7c201
todo create "Write auth tests" --type task --epic a3f7c201

todo list --epic a3f7c201
```

Create dependency chains:

```bash
todo create "Deploy to prod" --dep a3f7c201
# Won't appear in "todo ready" until a3f7c201 is closed
```

## Field Reference

| Field | Values | Default |
|---|---|---|
| type | task, bug, feature, epic, chore | task |
| status | open, in_progress, done, cancelled | open |
| priority | 0 (critical) to 4 (backlog) | 2 |
| epic | parent issue ID | — |
| deps | list of blocking issue IDs | — |
| labels | freeform string tags | — |

## Validation

Malformed issue files produce warnings with actionable details:

```
warning: bad00000-broken.md: validation errors:
  - invalid type "banana" (must be task, bug, feature, epic, or chore)
  - invalid status "pending" (must be open, in_progress, done, or cancelled)
  - invalid priority 9 (must be 0-4)
  - missing required field: created
```

## Claude Code Integration

A skill file is included at `skills/todo/SKILL.md` for use with [Claude Code](https://docs.anthropic.com/en/docs/claude-code). It instructs the LLM to use the CLI for queries and edit markdown files directly for writes.
