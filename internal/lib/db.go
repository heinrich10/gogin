package lib

import (
	"database/sql"
	"gogin/internal/config"
	"log/slog"

	_ "modernc.org/sqlite"
)

func GetConnection() (*sql.DB, error) {
	slog.Info("func", "GetConnection", slog.String("connecting", config.Config.DB_HOST))
	db, err := sql.Open("sqlite", config.Config.DB_HOST)
	if err != nil {
		slog.Error("func", "GetConnection", err)
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if err := db.Ping(); err != nil {
		slog.Error("func", "GetConnection", err)
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
