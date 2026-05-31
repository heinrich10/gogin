package lib

import (
	"database/sql"
	"errors"
	"fmt"
	"gogin/internal/config"
	"log/slog"

	_ "modernc.org/sqlite"
)

func GetConnection(cfg *config.Config) (*sql.DB, error) {
	if cfg == nil {
		return nil, errors.New("database config is nil")
	}
	if cfg.DB_HOST == "" {
		return nil, errors.New("DB_HOST is empty")
	}

	slog.Info("func", "GetConnection", slog.String("connecting", cfg.DB_HOST))
	db, err := sql.Open("sqlite", cfg.DB_HOST)
	if err != nil {
		slog.Error("func", "GetConnection", err)
		return nil, err
	}

	// Enable WAL mode
	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		slog.Error("func", "GetConnection", slog.String("error", "failed to enable WAL"), slog.Any("err", err))
		_ = db.Close()
		return nil, fmt.Errorf("failed to enable WAL: %w", err)
	}

	db.SetMaxOpenConns(10) // Increase connection pool
	if err := db.Ping(); err != nil {
		slog.Error("func", "GetConnection", err)
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
