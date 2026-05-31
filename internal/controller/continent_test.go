package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"gogin/internal/model"
	"gogin/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestContinentController_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(service.MockContinentService)
		ctrl := ContinentController{Service: mockService}

		continents := []model.Continent{
			{Code: "AF", Name: "Africa"},
			{Code: "AN", Name: "Antarctica"},
		}

		mockService.On("GetMany", mock.Anything, 10, 0).Return(continents, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/continents", nil)

		ctrl.Get(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []model.Continent
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, continents, response)
		mockService.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockService := new(service.MockContinentService)
		ctrl := ContinentController{Service: mockService}

		mockService.On("GetMany", mock.Anything, 10, 0).Return([]model.Continent{}, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/continents", nil)

		ctrl.Get(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestContinentController_GetOne(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(service.MockContinentService)
		ctrl := ContinentController{Service: mockService}

		continent := model.Continent{Code: "AF", Name: "Africa"}
		mockService.On("GetContinentByCode", mock.Anything, "AF").Return(continent, nil)

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
		mockService.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService := new(service.MockContinentService)
		ctrl := ContinentController{Service: mockService}

		mockService.On("GetContinentByCode", mock.Anything, "XX").Return(model.Continent{}, sql.ErrNoRows)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "code", Value: "XX"}}
		c.Request = httptest.NewRequest("GET", "/continents/XX", nil)

		ctrl.GetOne(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockService := new(service.MockContinentService)
		ctrl := ContinentController{Service: mockService}

		mockService.On("GetContinentByCode", mock.Anything, "AF").Return(model.Continent{}, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "code", Value: "AF"}}
		c.Request = httptest.NewRequest("GET", "/continents/AF", nil)

		ctrl.GetOne(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}
