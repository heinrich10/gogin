package lib

import (
	"fmt"
	"os"
	"path/filepath"
)

// MigrationsDir walks up from the current working directory until it finds
// go.mod and returns the migrations directory under the project root.
func MigrationsDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find project root containing go.mod")
		}
		dir = parent
	}
	return filepath.Join(dir, "migrations"), nil
}
