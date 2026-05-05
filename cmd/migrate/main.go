package main

import (
	"context"
	"gogin/internal/lib"
	"log/slog"
	"os"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

func migrationsDir(logger *slog.Logger) string {
	dir, err := lib.MigrationsDir()
	if err != nil {
		logger.Error("failed to locate migrations directory", slog.Any("error", err))
		os.Exit(1)
	}
	return dir
}

func main() {
	logger := slog.Default()

	db, err := lib.GetConnection()
	if err != nil {
		logger.Error("failed to get database connection", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	if err := goose.SetDialect("sqlite3"); err != nil {
		logger.Error("failed to set dialect", slog.Any("error", err))
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		logger.Error("missing command", slog.String("usage", "go run cmd/migrate/main.go <command> [args...]"))
		os.Exit(1)
	}

	command := os.Args[1]
	var args []string
	if len(os.Args) > 2 {
		args = append(args, os.Args[2:]...)
	}

	ctx := context.Background()
	if err := goose.RunContext(ctx, command, db, migrationsDir(logger), args...); err != nil {
		logger.Error("migration failed", slog.String("command", command), slog.Any("error", err))
		os.Exit(1)
	}
}
