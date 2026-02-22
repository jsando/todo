---
name: todo
description: This skill activates when a `.todo/` directory exists in the project, or when the user asks to create/manage issues with `todo`.
---

# todo — Markdown Issue Tracker Skill

## Activation

This skill activates when a `.todo/` directory exists in the project, or when the user asks to create/manage issues with `todo`.

## Overview

`todo` is a CLI tool that stores issues as individual markdown files with YAML front matter in `.todo/`. Each file is named `<8-char-hex-id>-<slug>.md`.

## For Queries — Use the CLI

```bash
todo list [--status X] [--type X] [--priority X] [--epic X] [--label X] [--json]
todo show <id> [--json]
todo ready [--json]          # issues with all deps resolved
todo blocked [--json]        # issues with unresolved deps
todo status [--json]         # summary counts
todo search <query> [--json]
```

Use `--json` for programmatic consumption.

## For Writes — Edit Files Directly

Create and update issues by editing markdown files directly with Read/Edit/Write tools. This is faster and more flexible than CLI commands.

### Issue Format

```markdown
---
id: a3f7c201
title: Add JWT authentication
type: feature
status: open
priority: 1
epic: c44b2e10
deps: [b912de44]
created: 2026-02-01
updated: 2026-02-01
labels: [auth, api]
---

Description goes here. Free-form markdown.
```

### Field Reference

- **id**: 8-char random hex
- **type**: task | bug | feature | epic | chore
- **status**: open | in_progress | done | cancelled
- **priority**: 0 (critical) to 4 (backlog), default 2
- **epic**: parent epic ID (optional)
- **deps**: list of blocking issue IDs (optional)
- **labels**: freeform string tags (optional)

### File Naming

Files are named `<id>-<slugified-title>.md`. When creating issues, generate an 8-char hex ID and slugify the title.

### Closing Issues

Move the file from `.todo/` to `.todo/archive/` and set `status: done`.

### Reopening Issues

Move from `.todo/archive/` back to `.todo/` and set `status: open`.

## Workflow

1. Check `todo ready --json` for unblocked work
2. Claim by setting `status: in_progress`
3. Do the work
4. Close with `todo close <id>` or move file to archive
5. Always commit `.todo/` changes alongside code changes

## Epics

Create an epic, then link child issues to it with `--epic`:

```bash
todo create "User auth system" --type epic --priority 1
# prints epic ID, e.g. a3f7c201
todo create "Design auth API" --type task --epic a3f7c201
todo create "Implement JWT tokens" --type task --epic a3f7c201
todo create "Write auth tests" --type task --epic a3f7c201

# List all issues in an epic:
todo list --epic a3f7c201

# List all epics:
todo list --type epic
```

## CLI Quick Reference

```bash
todo init                    # create .todo/ directory
todo create "title" --type X --priority X [--epic id] [--dep id]
todo update <id> --status X --add-label X --add-dep X
todo close <id> [--reason X] # move to archive
todo reopen <id>             # restore from archive
```

Partial IDs work (like git): `todo show a3f7` matches `a3f7c201-add-jwt-auth.md`.
