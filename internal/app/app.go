package app

import (
	"database/sql"
	"gogin/internal/controller"
	"gogin/internal/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewRouter builds a fully wired Gin engine with all controllers and routes.
// It also returns the UpdatePerson channel so the caller can manage its
// lifecycle (e.g. close it on shutdown).
func NewRouter(db *sql.DB) (*gin.Engine, chan controller.UpdatePerson) {
	continentRepository := repository.ContinentRepository{Db: db}
	countryRepository := repository.CountryRepository{Db: db}
	personRepository := repository.PersonRepository{Db: db}

	continentController := controller.ContinentController{
		Repository: &continentRepository,
	}

	countryController := controller.CountryController{
		Repository: &countryRepository,
	}

	updatePersonChan := make(chan controller.UpdatePerson, 100)
	personController := controller.PersonController{
		Repository:       &personRepository,
		UpdatePersonChan: updatePersonChan,
	}

	go personController.StartWorker()

	router := gin.Default()
	router.Use(cors.Default())

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

	return router, updatePersonChan
}
