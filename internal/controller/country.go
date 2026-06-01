package controller

import (
	"database/sql"
	"errors"
	"gogin/internal/service"
	"gogin/internal/util"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CountryController struct {
	Service service.CountryServiceInterface
}

func (d CountryController) Get(c *gin.Context) {
	slog.Info("func", "GetMany", slog.String("ip", c.ClientIP()))

	limit, offset, page := util.Paginate(c)

	rs, err := d.Service.GetMany(c.Request.Context(), limit, offset)
	if err != nil {
		slog.Error("failed to get countries", "ip", c.ClientIP(), "err", err)
		writeError(c, http.StatusInternalServerError, "ERR_INTERNAL", "Something went wrong")
		return
	}

	total, err := d.Service.Count(c.Request.Context())
	if err != nil {
		slog.Error("failed to count countries", "ip", c.ClientIP(), "err", err)
		writeError(c, http.StatusInternalServerError, "ERR_INTERNAL", "Something went wrong")
		return
	}

	c.JSON(http.StatusOK, NewListResponse(rs, page, limit, total))
}

func (d CountryController) GetOne(c *gin.Context) {
	slog.Info("func", "Get", slog.String("ip", c.ClientIP()))
	code := c.Param("code")
	rs, err := d.Service.GetCountryByCode(c.Request.Context(), code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(c, http.StatusNotFound, "ERR_NOT_FOUND", "Country not found")
			return
		}
		slog.Error("failed to get country by code", "code", code, "ip", c.ClientIP(), "err", err)
		writeError(c, http.StatusInternalServerError, "ERR_INTERNAL", "Something went wrong")
		return
	}
	c.JSON(http.StatusOK, rs)
}
