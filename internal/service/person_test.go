package service

import (
	"context"
	"errors"
	"gogin/internal/model"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPersonService_GetMany(t *testing.T) {
	mockRepo := new(MockPersonRepository)
	svc := &PersonService{Repo: mockRepo}

	ctx := context.Background()
	persons := []model.Person{{Id: 1, FirstName: "John"}}
	mockRepo.On("GetMany", ctx, 10, 0).Return(persons, nil)

	result, err := svc.GetMany(ctx, 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, persons, result)
	mockRepo.AssertExpectations(t)
}

func TestPersonService_GetPersonById(t *testing.T) {
	mockRepo := new(MockPersonRepository)
	svc := &PersonService{Repo: mockRepo}

	ctx := context.Background()
	person := model.Person{Id: 1, FirstName: "John"}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetPersonById", ctx, "1").Return(person, nil).Once()
		result, err := svc.GetPersonById(ctx, "1")
		assert.NoError(t, err)
		assert.Equal(t, person, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.On("GetPersonById", ctx, "99").Return(model.Person{}, errors.New("not found")).Once()
		result, err := svc.GetPersonById(ctx, "99")
		assert.Error(t, err)
		assert.Equal(t, model.Person{}, result)
	})

	mockRepo.AssertExpectations(t)
}

func TestPersonService_QueueCreate(t *testing.T) {
	updateChan := make(chan UpdatePerson, 1)
	svc := &PersonService{UpdatePersonChan: updateChan}
	person := model.Person{FirstName: "John"}

	svc.QueueCreate(person)

	select {
	case task := <-updateChan:
		assert.Equal(t, person, task.Person)
	case <-time.After(time.Millisecond * 100):
		t.Fatal("task was not queued")
	}
}

func TestPersonService_StartWorker(t *testing.T) {
	mockRepo := new(MockPersonRepository)
	updateChan := make(chan UpdatePerson, 5)

	// Use a cancelable context for the worker itself
	ctx, cancel := context.WithCancel(context.Background())
	// Use another for draining
	shutdownCtx, cancelShutdown := context.WithCancel(context.Background())
	defer cancelShutdown()

	svc := &PersonService{
		Repo:             mockRepo,
		UpdatePersonChan: updateChan,
		ShutdownCtx:      shutdownCtx,
	}

	person1 := model.Person{FirstName: "John"}
	person2 := model.Person{FirstName: "Jane"}

	mockRepo.On("Create", mock.Anything, person1).Return(nil)
	mockRepo.On("Create", mock.Anything, person2).Return(nil)

	var wg sync.WaitGroup
	wg.Add(1)
	go svc.StartWorker(ctx, &wg)

	// Send one task normally
	updateChan <- UpdatePerson{Person: person1}

	// Wait a bit to ensure it's processed
	time.Sleep(50 * time.Millisecond)

	// Send another task and cancel context immediately to test draining
	updateChan <- UpdatePerson{Person: person2}
	cancel()

	// Worker should drain person2 and finish
	wg.Wait()

	mockRepo.AssertExpectations(t)
}

func TestPersonService_StartWorker_CloseChannel(t *testing.T) {
	mockRepo := new(MockPersonRepository)
	updateChan := make(chan UpdatePerson, 1)
	svc := &PersonService{
		Repo:             mockRepo,
		UpdatePersonChan: updateChan,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go svc.StartWorker(context.Background(), &wg)

	close(updateChan)
	wg.Wait() // Should exit when channel is closed
}
