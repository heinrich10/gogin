package service

import (
	"context"
	"gogin/internal/model"
	"gogin/internal/repository"
	"log/slog"
	"sync"
)

type UpdatePerson struct {
	Person model.Person
}

type PersonServiceInterface interface {
	GetMany(ctx context.Context, limit, offset int) ([]model.Person, error)
	GetPersonById(ctx context.Context, id string) (model.Person, error)
	Count(ctx context.Context) (int, error)
	QueueCreate(person model.Person)
	StartWorker(ctx context.Context, wg *sync.WaitGroup)
}

type PersonService struct {
	Repo             repository.PersonRepositoryInterface
	UpdatePersonChan chan UpdatePerson
	ShutdownCtx      context.Context
}

func (s *PersonService) GetMany(ctx context.Context, limit, offset int) ([]model.Person, error) {
	return s.Repo.GetMany(ctx, limit, offset)
}

func (s *PersonService) GetPersonById(ctx context.Context, id string) (model.Person, error) {
	return s.Repo.GetPersonById(ctx, id)
}

func (s *PersonService) Count(ctx context.Context) (int, error) {
	return s.Repo.Count(ctx)
}

func (s *PersonService) QueueCreate(person model.Person) {
	s.UpdatePersonChan <- UpdatePerson{Person: person}
}

func (s *PersonService) StartWorker(ctx context.Context, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	slog.Info("func", "StartPersonWorker", "Starting person worker...")

	if ctx == nil {
		ctx = context.Background()
	}

	for {
		select {
		case <-ctx.Done():
			slog.Info("func", "StartWorker", "Context cancelled, draining channel...")
			drainCtx := s.ShutdownCtx
			if drainCtx == nil {
				drainCtx = context.Background()
			}
			for {
				select {
				case task, ok := <-s.UpdatePersonChan:
					if !ok {
						slog.Info("func", "StartWorker", "Worker finished draining.")
						return
					}
					s.processTask(drainCtx, task)
				default:
					slog.Info("func", "StartWorker", "Worker finished draining (no more tasks).")
					return
				}
			}
		case task, ok := <-s.UpdatePersonChan:
			if !ok {
				slog.Info("func", "StartWorker", "Channel closed, worker stopping.")
				return
			}
			s.processTask(ctx, task)
		}
	}
}

func (s *PersonService) processTask(ctx context.Context, task UpdatePerson) {
	slog.Info("func", "processTask", slog.String("processing", task.Person.FirstName))
	if err := s.Repo.Create(ctx, task.Person); err != nil {
		slog.Error("func", "processTask", err)
	} else {
		slog.Info("func", "processTask", slog.String("created", task.Person.FirstName))
	}
}
