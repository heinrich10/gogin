package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"gogin/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockContinentRepository is a mock for ContinentRepositoryInterface
type MockContinentRepository struct {
	mock.Mock
}

func (m *MockContinentRepository) GetContinentByCode(code string) (model.Continent, error) {
	args := m.Called(code)
	return args.Get(0).(model.Continent), args.Error(1)
}

func (m *MockContinentRepository) GetMany(limit, offset int) ([]model.Continent, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]model.Continent), args.Error(1)
}

func TestContinentController_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockContinentRepository)
		ctrl := ContinentController{Repository: mockRepo}

		continents := []model.Continent{
			{Code: "AF", Name: "Africa"},
			{Code: "AN", Name: "Antarctica"},
		}

		mockRepo.On("GetMany", 10, 0).Return(continents, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/continents", nil)

		ctrl.Get(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []model.Continent
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, continents, response)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := new(MockContinentRepository)
		ctrl := ContinentController{Repository: mockRepo}

		mockRepo.On("GetMany", 10, 0).Return([]model.Continent{}, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/continents", nil)

		ctrl.Get(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestContinentController_GetOne(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockContinentRepository)
		ctrl := ContinentController{Repository: mockRepo}

		continent := model.Continent{Code: "AF", Name: "Africa"}
		mockRepo.On("GetContinentByCode", "AF").Return(continent, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "code", Value: "AF"}}
		c.Request = httptest.NewRequest("GET", "/continents/AF", nil)

		ctrl.GetOne(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.Continent
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, continent, response)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo := new(MockContinentRepository)
		ctrl := ContinentController{Repository: mockRepo}

		mockRepo.On("GetContinentByCode", "XX").Return(model.Continent{}, sql.ErrNoRows)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "code", Value: "XX"}}
		c.Request = httptest.NewRequest("GET", "/continents/XX", nil)

		ctrl.GetOne(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := new(MockContinentRepository)
		ctrl := ContinentController{Repository: mockRepo}

		mockRepo.On("GetContinentByCode", "AF").Return(model.Continent{}, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "code", Value: "AF"}}
		c.Request = httptest.NewRequest("GET", "/continents/AF", nil)

		ctrl.GetOne(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockRepo.AssertExpectations(t)
	})
}
