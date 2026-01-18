package controller

import (
	"database/sql"
	"errors"
	"gogin/internal/repository"
	"net/http"
	"strconv"

	"log/slog"

	"github.com/gin-gonic/gin"
)

type ContinentController struct {
	Repository *repository.ContinentRepository
}

func (d ContinentController) Get(c *gin.Context) {
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

func (d ContinentController) GetOne(c *gin.Context) {
	slog.Info("func", "GetOne", slog.String("ip", c.ClientIP()))
	id := c.Param("code")
	rs, err := d.Repository.GetContinentByCode(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Continent not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, rs)
}
