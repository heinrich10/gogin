package controller

import (
	"context"
	"database/sql"
	"errors"
	"gogin/internal/model"
	"gogin/internal/repository"
	"gogin/internal/util"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type UpdatePerson struct {
	Person model.Person
}
type PersonController struct {
	Repository       repository.PersonRepositoryInterface
	UpdatePersonChan chan UpdatePerson
}

func (d PersonController) StartWorker(ctx context.Context, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	slog.Info("func", "StartPersonWorker", "Starting person worker...")

	if ctx == nil {
		ctx = context.Background()
	}

	for {
		select {
		case <-ctx.Done():
			slog.Info("func", "StartWorker", "Context cancelled, draining channel...")
			// Close the channel if we're draining to prevent further sends and ensure the loop terminates
			// Wait, we don't own the channel here usually, but if we are draining, we need to know when to stop.
			// The issue says: "On ctx.Done(), the worker drains by ranging over UpdatePersonChan, which blocks until the channel is closed."
			// We should probably check the channel in a non-blocking way or use a different approach.
			for {
				select {
				case task, ok := <-d.UpdatePersonChan:
					if !ok {
						slog.Info("func", "StartWorker", "Worker finished draining.")
						return
					}
					d.processTask(task)
				default:
					slog.Info("func", "StartWorker", "Worker finished draining (no more tasks).")
					return
				}
			}
		case task, ok := <-d.UpdatePersonChan:
			if !ok {
				slog.Info("func", "StartWorker", "Channel closed, worker stopping.")
				return
			}
			d.processTask(task)
		}
	}
}

func (d PersonController) processTask(task UpdatePerson) {
	slog.Info("func", "processTask", slog.String("processing", task.Person.FirstName))
	if err := d.Repository.Create(task.Person); err != nil {
		slog.Error("func", "processTask", err)
	} else {
		slog.Info("func", "processTask", slog.String("created", task.Person.FirstName))
	}
}

func (d PersonController) Get(c *gin.Context) {
	slog.Info("func", "GetMany", slog.String("ip", c.ClientIP()))

	limit, offset := util.Paginate(c)

	rs, err := d.Repository.GetMany(limit, offset)
	if err != nil {
		slog.Error("failed to get persons", "ip", c.ClientIP(), "err", err)
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
		slog.Error("failed to get person by id", "id", id, "ip", c.ClientIP(), "err", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	c.IndentedJSON(http.StatusOK, rs)
}

func (d PersonController) Create(c *gin.Context) {
	slog.Info("func", "Create", slog.String("ip", c.ClientIP()))
	var body model.Person
	if err := c.ShouldBindJSON(&body); err != nil {
		// Use a simple error message for now as we haven't implemented detailed validation error mapping yet,
		// but make sure it includes "validation failed" if it's a validation error.
		// Actually, let's just return the error string if it's not a syntax error.
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body or validation failed: " + err.Error()})
		return
	}
	d.UpdatePersonChan <- UpdatePerson{Person: body}
	c.IndentedJSON(http.StatusAccepted, gin.H{"status": "queued"})
}
