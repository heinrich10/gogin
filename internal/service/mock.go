package service

import (
	"context"
	"gogin/internal/model"
	"sync"

	"github.com/stretchr/testify/mock"
)

type MockContinentService struct {
	mock.Mock
}

func (m *MockContinentService) GetMany(ctx context.Context, limit, offset int) ([]model.Continent, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]model.Continent), args.Error(1)
}

func (m *MockContinentService) GetContinentByCode(ctx context.Context, code string) (model.Continent, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(model.Continent), args.Error(1)
}

type MockCountryService struct {
	mock.Mock
}

func (m *MockCountryService) GetMany(ctx context.Context, limit, offset int) ([]model.Country, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]model.Country), args.Error(1)
}

func (m *MockCountryService) GetCountryByCode(ctx context.Context, code string) (model.Country, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(model.Country), args.Error(1)
}

type MockPersonService struct {
	mock.Mock
}

func (m *MockPersonService) GetMany(ctx context.Context, limit, offset int) ([]model.Person, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]model.Person), args.Error(1)
}

func (m *MockPersonService) GetPersonById(ctx context.Context, id string) (model.Person, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Person), args.Error(1)
}

func (m *MockPersonService) QueueCreate(person model.Person) {
	m.Called(person)
}

func (m *MockPersonService) StartWorker(ctx context.Context, wg *sync.WaitGroup) {
	m.Called(ctx, wg)
}

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
