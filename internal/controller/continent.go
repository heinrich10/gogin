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

	limit, offset, page := util.Paginate(c)

	rs, err := d.Service.GetMany(c.Request.Context(), limit, offset)
	if err != nil {
		slog.Error("failed to get continents", "ip", c.ClientIP(), "err", err)
		writeError(c, http.StatusInternalServerError, "ERR_INTERNAL", "Something went wrong")
		return
	}

	total, err := d.Service.Count(c.Request.Context())
	if err != nil {
		slog.Error("failed to count continents", "ip", c.ClientIP(), "err", err)
		writeError(c, http.StatusInternalServerError, "ERR_INTERNAL", "Something went wrong")
		return
	}

	c.JSON(http.StatusOK, NewListResponse(rs, page, limit, total))
}

func (d ContinentController) GetOne(c *gin.Context) {
	slog.Info("func", "GetOne", slog.String("ip", c.ClientIP()))
	id := c.Param("code")
	rs, err := d.Service.GetContinentByCode(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(c, http.StatusNotFound, "ERR_NOT_FOUND", "Continent not found")
			return
		}
		slog.Error("failed to get continent by code", "code", id, "ip", c.ClientIP(), "err", err)
		writeError(c, http.StatusInternalServerError, "ERR_INTERNAL", "Something went wrong")
		return
	}
	c.JSON(http.StatusOK, rs)
}
