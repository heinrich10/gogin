package controller

import (
	"database/sql"
	"errors"
	"gogin/internal/model"
	"gogin/internal/service"
	"gogin/internal/util"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PersonController struct {
	Service service.PersonServiceInterface
}

func (d *PersonController) Get(c *gin.Context) {
	slog.Info("func", "GetMany", slog.String("ip", c.ClientIP()))

	limit, offset, page := util.Paginate(c)

	rs, err := d.Service.GetMany(c.Request.Context(), limit, offset)
	if err != nil {
		slog.Error("failed to get persons", "ip", c.ClientIP(), "err", err)
		writeError(c, http.StatusInternalServerError, "ERR_INTERNAL", "Something went wrong")
		return
	}

	total, err := d.Service.Count(c.Request.Context())
	if err != nil {
		slog.Error("failed to count persons", "ip", c.ClientIP(), "err", err)
		writeError(c, http.StatusInternalServerError, "ERR_INTERNAL", "Something went wrong")
		return
	}

	c.JSON(http.StatusOK, NewListResponse(rs, page, limit, total))
}

func (d *PersonController) GetOne(c *gin.Context) {
	slog.Info("func", "Get", slog.String("ip", c.ClientIP()))
	id := c.Param("id")
	rs, err := d.Service.GetPersonById(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(c, http.StatusNotFound, "ERR_NOT_FOUND", "Person not found")
			return
		}
		slog.Error("failed to get person by id", "id", id, "ip", c.ClientIP(), "err", err)
		writeError(c, http.StatusInternalServerError, "ERR_INTERNAL", "Something went wrong")
		return
	}
	c.JSON(http.StatusOK, rs)
}

func (d *PersonController) Create(c *gin.Context) {
	slog.Info("func", "Create", slog.String("ip", c.ClientIP()))
	var body model.Person
	if err := c.ShouldBindJSON(&body); err != nil {
		// Use a simple error message for now as we haven't implemented detailed validation error mapping yet,
		// but make sure it includes "validation failed" if it's a validation error.
		// Actually, let's just return the error string if it's not a syntax error.
		writeError(c, http.StatusBadRequest, "ERR_BAD_REQUEST", "Invalid request body or validation failed: "+err.Error())
		return
	}
	d.Service.QueueCreate(body)
	c.JSON(http.StatusAccepted, gin.H{"status": "queued"})
}
