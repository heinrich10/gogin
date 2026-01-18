package controller

import (
	"gogin/internal/repository"
	"net/http"
	"strconv"

	"log/slog"

	"github.com/gin-gonic/gin"
)

type CountryController struct {
	Repository *repository.CountryRepository
}

func (d CountryController) Get(c *gin.Context) {
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

func (d CountryController) GetOne(c *gin.Context) {
	slog.Info("func", "Get", slog.String("ip", c.ClientIP()))
	code := c.Param("code")
	rs, _ := d.Repository.GetCountryByCode(code)
	c.IndentedJSON(http.StatusOK, rs)
}
