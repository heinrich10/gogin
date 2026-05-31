package controller

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"gogin/internal/model"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPersonRepository is a mock for PersonRepositoryInterface
type MockPersonRepository struct {
	mock.Mock
}

func (m *MockPersonRepository) GetPersonById(ctx context.Context, id string) (model.Person, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Person), args.Error(1)
}

func (m *MockPersonRepository) GetMany(ctx context.Context, limit, offset int) ([]model.Person, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]model.Person), args.Error(1)
}

func (m *MockPersonRepository) Create(ctx context.Context, body model.Person) error {
	args := m.Called(ctx, body)
	return args.Error(0)
}

func TestPersonController_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockPersonRepository)
		ctrl := PersonController{Repository: mockRepo}

		persons := []model.Person{
			{Id: 1, FirstName: "John", LastName: "Doe"},
		}

		mockRepo.On("GetMany", mock.Anything, 10, 0).Return(persons, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/persons", nil)

		ctrl.Get(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []model.Person
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, persons, response)
		mockRepo.AssertExpectations(t)
	})
}

func TestPersonController_GetOne(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockPersonRepository)
		ctrl := PersonController{Repository: mockRepo}

		person := model.Person{Id: 1, FirstName: "John", LastName: "Doe"}
		mockRepo.On("GetPersonById", mock.Anything, "1").Return(person, nil)

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
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo := new(MockPersonRepository)
		ctrl := PersonController{Repository: mockRepo}

		mockRepo.On("GetPersonById", mock.Anything, "99").Return(model.Person{}, sql.ErrNoRows)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "99"}}
		c.Request = httptest.NewRequest("GET", "/persons/99", nil)

		ctrl.GetOne(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestPersonController_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		updateChan := make(chan UpdatePerson, 1)
		ctrl := PersonController{UpdatePersonChan: updateChan}

		person := model.Person{FirstName: "John", LastName: "Doe", CountryCode: "US"}
		body, _ := json.Marshal(person)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/persons", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		ctrl.Create(c)

		assert.Equal(t, http.StatusAccepted, w.Code)

		select {
		case task := <-updateChan:
			assert.Equal(t, person.FirstName, task.Person.FirstName)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Task not received on channel")
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		ctrl := PersonController{}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/persons", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		ctrl.Create(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
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
	})
}

func TestPersonController_StartWorker(t *testing.T) {
	mockRepo := new(MockPersonRepository)
	updateChan := make(chan UpdatePerson, 1) // Buffered to avoid blocking
	ctrl := PersonController{
		Repository:       mockRepo,
		UpdatePersonChan: updateChan,
	}

	person := model.Person{FirstName: "John", LastName: "Doe"}
	mockRepo.On("Create", mock.Anything, person).Return(nil)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go ctrl.StartWorker(ctx, &wg)

	updateChan <- UpdatePerson{Person: person}

	// Wait for worker to finish processing and draining
	close(updateChan)
	wg.Wait()
	cancel() // Cleanup context as well

	mockRepo.AssertExpectations(t)
}
