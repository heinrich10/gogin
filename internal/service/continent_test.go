package service

import (
	"context"
	"errors"
	"gogin/internal/model"
	"gogin/internal/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContinentService_GetMany(t *testing.T) {
	mockRepo := new(testutil.MockContinentRepository)
	svc := &ContinentService{Repo: mockRepo}

	ctx := context.Background()
	continents := []model.Continent{{Code: "EU", Name: "Europe"}}
	mockRepo.On("GetMany", ctx, 10, 0).Return(continents, nil)

	result, err := svc.GetMany(ctx, 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, continents, result)
	mockRepo.AssertExpectations(t)
}

func TestContinentService_GetContinentByCode(t *testing.T) {
	mockRepo := new(testutil.MockContinentRepository)
	svc := &ContinentService{Repo: mockRepo}

	ctx := context.Background()
	continent := model.Continent{Code: "EU", Name: "Europe"}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetContinentByCode", ctx, "EU").Return(continent, nil).Once()
		result, err := svc.GetContinentByCode(ctx, "EU")
		assert.NoError(t, err)
		assert.Equal(t, continent, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.On("GetContinentByCode", ctx, "XX").Return(model.Continent{}, errors.New("not found")).Once()
		result, err := svc.GetContinentByCode(ctx, "XX")
		assert.Error(t, err)
		assert.Equal(t, model.Continent{}, result)
	})

	mockRepo.AssertExpectations(t)
}
