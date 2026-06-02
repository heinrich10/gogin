package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"gogin/internal/model"
	"gogin/internal/testutil"
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
		mockService := new(testutil.MockContinentService)
		ctrl := ContinentController{Service: mockService}

		continents := []model.Continent{
			{Code: "AF", Name: "Africa"},
			{Code: "AN", Name: "Antarctica"},
		}

		mockService.On("GetMany", mock.Anything, 10, 0).Return(continents, nil)
		mockService.On("Count", mock.Anything).Return(7, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/continents", nil)

		ctrl.Get(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response ListResponse[model.Continent]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, continents, response.Data)
		assert.Equal(t, 1, response.Meta.Page)
		assert.Equal(t, 10, response.Meta.Limit)
		assert.Equal(t, 7, response.Meta.Total)
		mockService.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockService := new(testutil.MockContinentService)
		ctrl := ContinentController{Service: mockService}

		mockService.On("GetMany", mock.Anything, 10, 0).Return([]model.Continent{}, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/continents", nil)

		ctrl.Get(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ERR_INTERNAL", response.Error.Code)
		mockService.AssertExpectations(t)
	})
}

func TestContinentController_GetOne(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(testutil.MockContinentService)
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
		mockService := new(testutil.MockContinentService)
		ctrl := ContinentController{Service: mockService}

		mockService.On("GetContinentByCode", mock.Anything, "XX").Return(model.Continent{}, sql.ErrNoRows)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "code", Value: "XX"}}
		c.Request = httptest.NewRequest("GET", "/continents/XX", nil)

		ctrl.GetOne(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ERR_NOT_FOUND", response.Error.Code)
		assert.Equal(t, "Continent not found", response.Error.Message)
		mockService.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockService := new(testutil.MockContinentService)
		ctrl := ContinentController{Service: mockService}

		mockService.On("GetContinentByCode", mock.Anything, "AF").Return(model.Continent{}, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "code", Value: "AF"}}
		c.Request = httptest.NewRequest("GET", "/continents/AF", nil)

		ctrl.GetOne(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ERR_INTERNAL", response.Error.Code)
		mockService.AssertExpectations(t)
	})
}
