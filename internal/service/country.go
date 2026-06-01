package service

import (
	"context"
	"gogin/internal/model"
	"gogin/internal/repository"
)

type CountryServiceInterface interface {
	GetMany(ctx context.Context, limit, offset int) ([]model.Country, error)
	GetCountryByCode(ctx context.Context, code string) (model.Country, error)
	Count(ctx context.Context) (int, error)
}

type CountryService struct {
	Repo repository.CountryRepositoryInterface
}

func (s *CountryService) GetMany(ctx context.Context, limit, offset int) ([]model.Country, error) {
	return s.Repo.GetMany(ctx, limit, offset)
}

func (s *CountryService) GetCountryByCode(ctx context.Context, code string) (model.Country, error) {
	return s.Repo.GetCountryByCode(ctx, code)
}

func (s *CountryService) Count(ctx context.Context) (int, error) {
	return s.Repo.Count(ctx)
}
