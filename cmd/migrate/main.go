package main

import (
	"fmt"
	"gogin/internal/config"
	"gogin/internal/lib"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	cfg := config.LoadConfig()

	db, err := lib.GetConnection()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Resolve migrations directory relative to module root.
	// When run with `go run ./cmd/migrate/main.go` the working directory is the module root.
	migrationFiles, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to list migrations: %v\n", err)
		os.Exit(1)
	}
	sort.Strings(migrationFiles)

	for _, f := range migrationFiles {
		content, err := os.ReadFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read %s: %v\n", f, err)
			os.Exit(1)
		}

		s := string(content)
		upStart := strings.Index(s, "-- migrate:up")
		downStart := strings.Index(s, "-- migrate:down")

		var upSQL string
		if upStart != -1 && downStart != -1 {
			upSQL = s[upStart+len("-- migrate:up") : downStart]
		} else if upStart != -1 {
			upSQL = s[upStart+len("-- migrate:up"):]
		}

		if _, err := db.Exec(upSQL); err != nil {
			fmt.Fprintf(os.Stderr, "failed to execute %s: %v\n", f, err)
			os.Exit(1)
		}
		fmt.Printf("applied: %s\n", filepath.Base(f))
	}

	fmt.Printf("migrations applied successfully to %s\n", cfg.DB_HOST)
}
