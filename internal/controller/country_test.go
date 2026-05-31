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

func TestCountryController_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(service.MockCountryService)
		ctrl := CountryController{Service: mockService}

		countries := []model.Country{
			{Code: "US", Name: "United States"},
			{Code: "CA", Name: "Canada"},
		}

		mockService.On("GetMany", mock.Anything, 10, 0).Return(countries, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/countries", nil)

		ctrl.Get(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []model.Country
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, countries, response)
		mockService.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockService := new(service.MockCountryService)
		ctrl := CountryController{Service: mockService}

		mockService.On("GetMany", mock.Anything, 10, 0).Return([]model.Country{}, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/countries", nil)

		ctrl.Get(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestCountryController_GetOne(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(service.MockCountryService)
		ctrl := CountryController{Service: mockService}

		country := model.Country{Code: "US", Name: "United States"}
		mockService.On("GetCountryByCode", mock.Anything, "US").Return(country, nil)

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
		mockService.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService := new(service.MockCountryService)
		ctrl := CountryController{Service: mockService}

		mockService.On("GetCountryByCode", mock.Anything, "XX").Return(model.Country{}, sql.ErrNoRows)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "code", Value: "XX"}}
		c.Request = httptest.NewRequest("GET", "/countries/XX", nil)

		ctrl.GetOne(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockService := new(service.MockCountryService)
		ctrl := CountryController{Service: mockService}

		mockService.On("GetCountryByCode", mock.Anything, "US").Return(model.Country{}, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "code", Value: "US"}}
		c.Request = httptest.NewRequest("GET", "/countries/US", nil)

		ctrl.GetOne(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}
