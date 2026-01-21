package controller

import (
	"database/sql"
	"errors"
	"gogin/internal/repository"
	"gogin/internal/util"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CountryController struct {
	Repository repository.CountryRepositoryInterface
}

func (d CountryController) Get(c *gin.Context) {
	slog.Info("func", "GetMany", slog.String("ip", c.ClientIP()))

	limit, offset := util.Paginate(c)

	rs, err := d.Repository.GetMany(limit, offset)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.IndentedJSON(http.StatusOK, rs)
}

func (d CountryController) GetOne(c *gin.Context) {
	slog.Info("func", "Get", slog.String("ip", c.ClientIP()))
	code := c.Param("code")
	rs, err := d.Repository.GetCountryByCode(code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Country not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	c.IndentedJSON(http.StatusOK, rs)
}
