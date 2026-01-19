package controller

import (
	"database/sql"
	"errors"
	"gogin/internal/model"
	"gogin/internal/repository"
	"gogin/internal/util"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdatePerson struct {
	Person model.Person
}
type PersonController struct {
	Repository       *repository.PersonRepository
	UpdatePersonChan chan UpdatePerson
}

func (d PersonController) StartWorker() {
	slog.Info("func", "StartPersonWorker", "Starting person worker...")
	for task := range d.UpdatePersonChan {
		slog.Info("func", "StartWorker", slog.String("processing", task.Person.FirstName))
		if err := d.Repository.Create(task.Person); err != nil {
			slog.Error("func", "StartWorker", err)
		} else {
			slog.Info("func", "StartWorker", slog.String("created", task.Person.FirstName))
		}
	}
}

func (d PersonController) Get(c *gin.Context) {
	slog.Info("func", "GetMany", slog.String("ip", c.ClientIP()))

	limit, offset := util.Paginate(c)

	rs, err := d.Repository.GetMany(limit, offset)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.IndentedJSON(http.StatusOK, rs)
}

func (d PersonController) GetOne(c *gin.Context) {
	slog.Info("func", "Get", slog.String("ip", c.ClientIP()))
	id := c.Param("id")
	rs, err := d.Repository.GetPersonById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Person not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	c.IndentedJSON(http.StatusOK, rs)
}

func (d PersonController) Create(c *gin.Context) {
	slog.Info("func", "Create", slog.String("ip", c.ClientIP()))
	var body model.Person
	if err := c.ShouldBindJSON(&body); err != nil {
		return
	}
	d.UpdatePersonChan <- UpdatePerson{Person: body}
	c.IndentedJSON(http.StatusOK, gin.H{"status": 1})
}
