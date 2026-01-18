package controller

import (
	"gogin/internal/model"
	"gogin/internal/repository"
	"net/http"
	"strconv"

	"log/slog"

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
		_ = d.Repository.Create(task.Person)
	}
}

func (d PersonController) Get(c *gin.Context) {
	slog.Info("func", "GetMany", slog.String("ip", c.ClientIP()))

	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	rs, err := d.Repository.GetMany(limit, offset)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, rs)
}

func (d PersonController) GetOne(c *gin.Context) {
	slog.Info("func", "Get", slog.String("ip", c.ClientIP()))
	id := c.Param("id")
	rs, _ := d.Repository.GetPersonById(id)
	c.IndentedJSON(http.StatusOK, rs)
}

func (d PersonController) Create(c *gin.Context) {
	slog.Info("func", "Create", slog.String("ip", c.ClientIP()))
	var body model.Person
	err := c.ShouldBindJSON(&body)
	if err != nil {
		return
	}
	d.UpdatePersonChan <- UpdatePerson{Person: body}
	c.IndentedJSON(http.StatusOK, gin.H{"status": 1})
}
