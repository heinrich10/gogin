package app

import (
	"context"
	"database/sql"
	"gogin/internal/config"
	"gogin/internal/controller"
	"gogin/internal/repository"
	"gogin/internal/service"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewRouter builds a fully wired Gin engine with all controllers and routes.
// It also returns the UpdatePerson channel so the caller can manage its
// lifecycle (e.g. close it on shutdown).
func NewRouter(ctx context.Context, shutdownCtx context.Context, wg *sync.WaitGroup, db *sql.DB, cfg *config.Config) (*gin.Engine, chan service.UpdatePerson) {
	continentRepository := repository.ContinentRepository{Db: db}
	countryRepository := repository.CountryRepository{Db: db}
	personRepository := repository.PersonRepository{Db: db}

	continentService := service.ContinentService{
		Repo: &continentRepository,
	}
	countryService := service.CountryService{
		Repo: &countryRepository,
	}

	updatePersonChan := make(chan service.UpdatePerson, 100)
	personService := service.PersonService{
		Repo:             &personRepository,
		UpdatePersonChan: updatePersonChan,
		ShutdownCtx:      shutdownCtx,
	}

	continentController := controller.ContinentController{
		Service: &continentService,
	}

	countryController := controller.CountryController{
		Service: &countryService,
	}

	personController := controller.PersonController{
		Service: &personService,
	}

	if wg != nil {
		wg.Add(1)
	}
	go personService.StartWorker(ctx, wg)

	allowCredentials := true
	for _, origin := range cfg.ALLOWED_ORIGINS {
		if origin == "*" {
			allowCredentials = false
			break
		}
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.ALLOWED_ORIGINS,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: allowCredentials,
		MaxAge:           12 * time.Hour,
	}))

	{
		continentGroup := router.Group("/continents")
		continentGroup.GET("", continentController.Get)
		continentGroup.GET("/:code", continentController.GetOne)
	}
	{
		countryGroup := router.Group("/countries")
		countryGroup.GET("", countryController.Get)
		countryGroup.GET("/:code", countryController.GetOne)
	}
	{
		personGroup := router.Group("/persons")
		personGroup.GET("", personController.Get)
		personGroup.GET("/:id", personController.GetOne)
		personGroup.POST("", personController.Create)
	}

	return router, updatePersonChan
}
