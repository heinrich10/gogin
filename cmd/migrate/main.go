package main

import (
	"context"
	"database/sql"
	"gogin/internal/config"
	"log"
	"os"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

func main() {
	cfg := config.LoadConfig()
	if cfg.DB_HOST == "" {
		log.Fatal("DB_HOST is not set")
	}

	db, err := sql.Open("sqlite", cfg.DB_HOST)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatalf("failed to set dialect: %v", err)
	}

	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <command> [args...]", os.Args[0])
	}

	command := os.Args[1]
	var args []string
	if len(os.Args) > 2 {
		args = append(args, os.Args[2:]...)
	}

	ctx := context.Background()
	if err := goose.RunContext(ctx, command, db, "migrations", args...); err != nil {
		log.Fatalf("goose %s: %v", command, err)
	}
}
