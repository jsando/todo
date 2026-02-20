#!/bin/sh
set -e

echo "Installing todo CLI..."
go install github.com/jsando/todo@latest

echo "Installing Claude Code skill..."
mkdir -p ~/.claude/skills/todo
curl -sL https://raw.githubusercontent.com/jsando/todo/main/skills/todo/SKILL.md \
  -o ~/.claude/skills/todo/SKILL.md

echo "Done. Run 'todo help' to get started."
