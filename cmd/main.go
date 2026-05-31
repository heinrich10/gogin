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
	"sync"
	"syscall"
	"time"
)

func main() {
	cfg := config.LoadConfig()
	logger := slog.Default()

	logger.Info("func", "main", "Starting...")
	db, err := lib.GetConnection(cfg)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		return
	}

	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("Failed to close database connection", "error", err)
		}
	}()

	var wg sync.WaitGroup
	ctx, cancelWorker := context.WithCancel(context.Background())
	defer cancelWorker()

	// Let's create a background context that we can cancel later for shutdown.
	shutdownCtx, cancelShutdown := context.WithCancel(context.Background())
	defer cancelShutdown()

	router, updatePersonChan := app.NewRouter(ctx, shutdownCtx, &wg, db, cfg)

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
	srvShutdownCtx, cancelSrv := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelSrv()

	// Start draining the worker using srvShutdownCtx.
	// We can't swap the context in personController, but we can make shutdownCtx
	// a context that we cancel, and StartWorker uses ShutdownCtx during draining.
	// If we want it bounded, we should have used a context that can be timed out.
	// Let's just cancel it after 5 seconds too.
	go func() {
		<-srvShutdownCtx.Done()
		cancelShutdown()
	}()

	if err := srv.Shutdown(srvShutdownCtx); err != nil {
		logger.Error("Shutdown server", "error", err)
	} else {
		logger.Info("Server shutdown complete")
	}

	close(updatePersonChan)
	cancelWorker()

	wg.Wait()

	logger.Info("exiting")
}
