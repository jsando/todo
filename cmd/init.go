package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

func Init(args []string) error {
	todoDir := ".todo"
	archiveDir := filepath.Join(todoDir, "archive")
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return err
	}
	fmt.Println("Initialized .todo/ directory")
	return nil
}
