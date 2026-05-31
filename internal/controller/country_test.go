package controller

import (
	"context"
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

// MockCountryRepository is a mock for CountryRepositoryInterface
type MockCountryRepository struct {
	mock.Mock
}

func (m *MockCountryRepository) GetCountryByCode(ctx context.Context, code string) (model.Country, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(model.Country), args.Error(1)
}

func (m *MockCountryRepository) GetMany(ctx context.Context, limit, offset int) ([]model.Country, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]model.Country), args.Error(1)
}

func TestCountryController_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockCountryRepository)
		ctrl := CountryController{Repository: mockRepo}

		countries := []model.Country{
			{Code: "US", Name: "United States"},
			{Code: "CA", Name: "Canada"},
		}

		mockRepo.On("GetMany", mock.Anything, 10, 0).Return(countries, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/countries", nil)

		ctrl.Get(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []model.Country
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, countries, response)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := new(MockCountryRepository)
		ctrl := CountryController{Repository: mockRepo}

		mockRepo.On("GetMany", mock.Anything, 10, 0).Return([]model.Country{}, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/countries", nil)

		ctrl.Get(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestCountryController_GetOne(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockCountryRepository)
		ctrl := CountryController{Repository: mockRepo}

		country := model.Country{Code: "US", Name: "United States"}
		mockRepo.On("GetCountryByCode", mock.Anything, "US").Return(country, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "code", Value: "US"}}
		c.Request = httptest.NewRequest("GET", "/countries/US", nil)

		ctrl.GetOne(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.Country
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, country, response)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo := new(MockCountryRepository)
		ctrl := CountryController{Repository: mockRepo}

		mockRepo.On("GetCountryByCode", mock.Anything, "XX").Return(model.Country{}, sql.ErrNoRows)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "code", Value: "XX"}}
		c.Request = httptest.NewRequest("GET", "/countries/XX", nil)

		ctrl.GetOne(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := new(MockCountryRepository)
		ctrl := CountryController{Repository: mockRepo}

		mockRepo.On("GetCountryByCode", mock.Anything, "US").Return(model.Country{}, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "code", Value: "US"}}
		c.Request = httptest.NewRequest("GET", "/countries/US", nil)

		ctrl.GetOne(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockRepo.AssertExpectations(t)
	})
}
