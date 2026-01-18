package main

import (
	"context"
	"fmt"
	"gogin/internal/config"
	"gogin/internal/controller"
	"gogin/internal/lib"
	"gogin/internal/repository"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	logger := slog.Default()

	logger.Info("func", "main", "Starting...")
	db, err := lib.GetConnection()
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		return
	}

	continentRepository := repository.ContinentRepository{Db: db}
	countryRepository := repository.CountryRepository{Db: db}
	personRepository := repository.PersonRepository{Db: db}

	updatePersonChan := make(chan controller.UpdatePerson, 100)

	continentController := controller.ContinentController{
		Repository: &continentRepository,
	}

	countryController := controller.CountryController{
		Repository: &countryRepository,
	}

	personController := controller.PersonController{
		Repository:       &personRepository,
		UpdatePersonChan: updatePersonChan,
	}

	go personController.StartWorker()

	router := gin.Default()
	router.Use(cors.Default())
	err = router.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		logger.Error("Failed to set trusted proxies", "error", err)
		return
	}
	{
		continentGroup := router.Group("/continents")
		continentGroup.GET("/", continentController.Get)
		continentGroup.GET("/:code", continentController.GetOne)
	}
	{
		countryGroup := router.Group("/countries")
		countryGroup.GET("/", countryController.Get)
		countryGroup.GET("/:code", countryController.GetOne)
	}
	{
		personGroup := router.Group("/persons")
		personGroup.GET("/", personController.Get)
		personGroup.GET("/:id", personController.GetOne)
		personGroup.POST("/", personController.Create)
	}

	hostPort := fmt.Sprintf("%s:%d", config.Config.HOST, config.Config.PORT)
	srv := &http.Server{
		Addr:    hostPort,
		Handler: router,
	}

	go func() {
		logger.Info(
			fmt.Sprintf(
				"Server running on %s",
				hostPort,
			),
		)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	close(updatePersonChan)
	if err := db.Close(); err != nil {
		logger.Error("Failed to close database connection", "error", err)
	}
	logger.Info("Server exiting")
}
