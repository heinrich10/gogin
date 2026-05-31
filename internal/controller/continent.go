package controller

import (
	"database/sql"
	"errors"
	"gogin/internal/service"
	"gogin/internal/util"
	"net/http"

	"log/slog"

	"github.com/gin-gonic/gin"
)

type ContinentController struct {
	Service service.ContinentServiceInterface
}

func (d ContinentController) Get(c *gin.Context) {
	slog.Info("func", "GetMany", slog.String("ip", c.ClientIP()))

	limit, offset := util.Paginate(c)

	rs, err := d.Service.GetMany(c.Request.Context(), limit, offset)
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
	rs, err := d.Service.GetContinentByCode(c.Request.Context(), id)
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
