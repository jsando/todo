# `todo` — Markdown-based Issue Tracker

## Overview

A lightweight CLI tool (Go) and Claude Code skill for project-level issue tracking using markdown files in a `.todo/` directory. Each issue is a separate `.md` file with YAML front matter. No central index or database — all queries scan the directory.

## File Layout

```
.todo/
  a3f7c201-add-jwt-auth.md
  b912de44-fix-login-bug.md
  archive/
    a3f7c201-add-jwt-auth.md    # completed/cancelled issues moved here
```

## Issue Format

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

- **id**: 8-char random hex (`crypto/rand`, 4 bytes)
- **type**: task | bug | feature | epic | chore
- **status**: open | in_progress | done | cancelled
- **priority**: 0 (critical) to 4 (backlog), default 2
- **epic**: id of parent epic (optional)
- **deps**: list of issue IDs this is blocked by
- **labels**: freeform string tags

## CLI Commands

Built as a single Go binary `todo`.

| Command | Description |
|---|---|
| `todo init` | Create `.todo/` and `.todo/archive/` |
| `todo create "title" [--type X] [--priority X] [--epic id] [--dep id]` | Create issue, print id |
| `todo list [--status X] [--type X] [--priority X] [--epic X] [--label X]` | List issues as table |
| `todo show <id>` | Print full issue file |
| `todo update <id> [--status X] [--priority X] [--title X] [--add-dep X] [--rm-dep X] [--add-label X] [--rm-label X] [--epic X]` | Update fields |
| `todo close <id> [--reason X]` | Set status=done, move to archive/ |
| `todo reopen <id>` | Move from archive/ back, set status=open |
| `todo ready` | List open issues with all deps resolved |
| `todo blocked` | List open issues with unresolved deps |
| `todo status` | Summary counts by status/type/priority |
| `todo search <query>` | Grep across all issue files |

All commands accept `--json` for LLM consumption.

### ID Resolution

Like git, partial IDs work. `todo show a3f7` matches `a3f7c201-add-jwt-auth.md`. Errors on ambiguous or no match.

## Go Project Structure

```
go.mod
main.go              # CLI entry point (subcommand dispatch with os.Args and flag.FlagSet)
cmd/
  args.go            # reorderArgs helper for mixed flag/positional parsing
  init.go
  create.go
  list.go
  show.go
  update.go
  close.go
  reopen.go
  ready.go
  blocked.go
  status.go
  search.go
internal/
  issue.go            # Issue struct, Parse(), Serialize() with YAML front matter
  store.go            # FindRoot(), LoadAll(), LoadArchived(), Resolve()
  format.go           # Table and JSON output formatting
skills/
  todo/
    SKILL.md          # Claude Code skill for LLM-assisted issue management
```

Dependencies: `gopkg.in/yaml.v3` for YAML. Standard library for CLI — manual subcommand dispatch with `os.Args` and `flag.FlagSet` per command.

## Claude Code Skill

The skill prompt instructs the LLM to:
- Detect `.todo/` directory to activate
- **Create/update issues**: Edit markdown files directly with Read/Edit/Write tools
- **Query issues**: Use `todo ready`, `todo list`, `todo blocked`, etc.
- Follow the YAML front matter schema exactly
- Always commit `.todo/` changes with code changes
