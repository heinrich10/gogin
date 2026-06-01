package controller

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"gogin/internal/model"
	"gogin/internal/service"
	"gogin/internal/testutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPersonController_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(testutil.MockPersonService)
		ctrl := PersonController{Service: mockService}

		persons := []model.Person{
			{Id: 1, FirstName: "John", LastName: "Doe"},
		}

		mockService.On("GetMany", mock.Anything, 10, 0).Return(persons, nil)
		mockService.On("Count", mock.Anything).Return(42, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/persons", nil)

		ctrl.Get(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response ListResponse[model.Person]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, persons, response.Data)
		assert.Equal(t, 1, response.Meta.Page)
		assert.Equal(t, 10, response.Meta.Limit)
		assert.Equal(t, 42, response.Meta.Total)
		mockService.AssertExpectations(t)
	})
}

func TestPersonController_GetOne(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(testutil.MockPersonService)
		ctrl := PersonController{Service: mockService}

		person := model.Person{Id: 1, FirstName: "John", LastName: "Doe"}
		mockService.On("GetPersonById", mock.Anything, "1").Return(person, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest("GET", "/persons/1", nil)

		ctrl.GetOne(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.Person
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, person, response)
		mockService.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService := new(testutil.MockPersonService)
		ctrl := PersonController{Service: mockService}

		mockService.On("GetPersonById", mock.Anything, "99").Return(model.Person{}, sql.ErrNoRows)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "99"}}
		c.Request = httptest.NewRequest("GET", "/persons/99", nil)

		ctrl.GetOne(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ERR_NOT_FOUND", response.Error.Code)
		assert.Equal(t, "Person not found", response.Error.Message)
		mockService.AssertExpectations(t)
	})
}

func TestPersonController_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(testutil.MockPersonService)
		ctrl := PersonController{Service: mockService}

		person := model.Person{FirstName: "John", LastName: "Doe", CountryCode: "US"}
		body, _ := json.Marshal(person)

		mockService.On("QueueCreate", person).Return()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/persons", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		ctrl.Create(c)

		assert.Equal(t, http.StatusAccepted, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		ctrl := PersonController{}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/persons", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		ctrl.Create(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ERR_BAD_REQUEST", response.Error.Code)
		assert.Contains(t, response.Error.Message, "Invalid request body")
	})

	t.Run("Validation Failure - Missing FirstName", func(t *testing.T) {
		ctrl := PersonController{}

		person := model.Person{LastName: "Doe", CountryCode: "US"}
		body, _ := json.Marshal(person)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/persons", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		ctrl.Create(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ERR_BAD_REQUEST", response.Error.Code)
		assert.Contains(t, response.Error.Message, "validation failed")
	})

	t.Run("Validation Failure - Invalid CountryCode", func(t *testing.T) {
		ctrl := PersonController{}

		person := model.Person{FirstName: "John", LastName: "Doe", CountryCode: "USA"}
		body, _ := json.Marshal(person)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/persons", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		ctrl.Create(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ERR_BAD_REQUEST", response.Error.Code)
		assert.Contains(t, response.Error.Message, "validation failed")
	})
}

func TestPersonController_StartWorker(t *testing.T) {
	mockRepo := new(testutil.MockPersonRepository)
	updateChan := make(chan service.UpdatePerson, 1) // Buffered to avoid blocking
	personService := service.PersonService{
		Repo:             mockRepo,
		UpdatePersonChan: updateChan,
	}

	person := model.Person{FirstName: "John", LastName: "Doe"}
	mockRepo.On("Create", mock.Anything, person).Return(nil)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go personService.StartWorker(ctx, &wg)

	updateChan <- service.UpdatePerson{Person: person}

	// Wait for worker to finish processing and draining
	close(updateChan)
	wg.Wait()
	cancel() // Cleanup context as well

	mockRepo.AssertExpectations(t)
}
