package controller

import (
	"database/sql"
	"errors"
	"gogin/internal/repository"
	"gogin/internal/util"
	"net/http"

	"log/slog"

	"github.com/gin-gonic/gin"
)

type ContinentController struct {
	Repository repository.ContinentRepositoryInterface
}

func (d ContinentController) Get(c *gin.Context) {
	slog.Info("func", "GetMany", slog.String("ip", c.ClientIP()))

	limit, offset := util.Paginate(c)

	rs, err := d.Repository.GetMany(limit, offset)
	if err != nil {
		slog.Error("failed to get continents", "ip", c.ClientIP(), "err", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
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
		slog.Error("failed to get continent by code", "code", id, "ip", c.ClientIP(), "err", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	c.IndentedJSON(http.StatusOK, rs)
}
