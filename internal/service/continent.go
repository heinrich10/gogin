package service

import (
	"context"
	"gogin/internal/model"
	"gogin/internal/repository"
)

type ContinentServiceInterface interface {
	GetMany(ctx context.Context, limit, offset int) ([]model.Continent, error)
	GetContinentByCode(ctx context.Context, code string) (model.Continent, error)
	Count(ctx context.Context) (int, error)
}

type ContinentService struct {
	Repo repository.ContinentRepositoryInterface
}

func (s *ContinentService) GetMany(ctx context.Context, limit, offset int) ([]model.Continent, error) {
	return s.Repo.GetMany(ctx, limit, offset)
}

func (s *ContinentService) GetContinentByCode(ctx context.Context, code string) (model.Continent, error) {
	return s.Repo.GetContinentByCode(ctx, code)
}

func (s *ContinentService) Count(ctx context.Context) (int, error) {
	return s.Repo.Count(ctx)
}
