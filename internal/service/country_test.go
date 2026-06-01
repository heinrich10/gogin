package service

import (
	"context"
	"errors"
	"gogin/internal/model"
	"gogin/internal/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountryService_GetMany(t *testing.T) {
	mockRepo := new(testutil.MockCountryRepository)
	svc := &CountryService{Repo: mockRepo}

	ctx := context.Background()
	countries := []model.Country{{Code: "US", Name: "United States"}}
	mockRepo.On("GetMany", ctx, 10, 0).Return(countries, nil)

	result, err := svc.GetMany(ctx, 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, countries, result)
	mockRepo.AssertExpectations(t)
}

func TestCountryService_GetCountryByCode(t *testing.T) {
	mockRepo := new(testutil.MockCountryRepository)
	svc := &CountryService{Repo: mockRepo}

	ctx := context.Background()
	country := model.Country{Code: "US", Name: "United States"}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetCountryByCode", ctx, "US").Return(country, nil).Once()
		result, err := svc.GetCountryByCode(ctx, "US")
		assert.NoError(t, err)
		assert.Equal(t, country, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.On("GetCountryByCode", ctx, "XX").Return(model.Country{}, errors.New("not found")).Once()
		result, err := svc.GetCountryByCode(ctx, "XX")
		assert.Error(t, err)
		assert.Equal(t, model.Country{}, result)
	})

	mockRepo.AssertExpectations(t)
}

func TestCountryService_Count(t *testing.T) {
	mockRepo := new(testutil.MockCountryRepository)
	svc := &CountryService{Repo: mockRepo}

	ctx := context.Background()
	mockRepo.On("Count", ctx).Return(250, nil)

	result, err := svc.Count(ctx)

	assert.NoError(t, err)
	assert.Equal(t, 250, result)
	mockRepo.AssertExpectations(t)
}
