package main

import (
	"context"
	"errors"
	"fmt"
	"gogin/internal/app"
	"gogin/internal/config"
	"gogin/internal/lib"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.LoadConfig()
	logger := slog.Default()

	logger.Info("func", "main", "Starting...")
	db, err := lib.GetConnection()
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		return
	}

	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("Failed to close database connection", "error", err)
		}
	}()

	router, updatePersonChan := app.NewRouter(db)

	if err := router.SetTrustedProxies(cfg.TRUSTED_PROXIES); err != nil {
		logger.Error("Failed to set trusted proxies", "error", err)
		return
	}

	hostPort := fmt.Sprintf("%s:%d", cfg.HOST, cfg.PORT)
	srv := &http.Server{
		Addr:    hostPort,
		Handler: router,
	}

	go func() {
		logger.Info("Server running on", "address", hostPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Failed to start server", "error", err)
		}
	}()

	// Wait for termination signal.
	sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-sigCtx.Done()
	logger.Info("Shutdown requested")

	// Graceful shutdown with timeout.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Shutdown server", "error", err)
	} else {
		logger.Info("Server shutdown complete")
	}

	close(updatePersonChan)

	// small grace period for worker cleanup (adjust if needed)
	time.Sleep(100 * time.Millisecond)

	logger.Info("exiting")
}
