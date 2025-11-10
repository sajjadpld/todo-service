package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	localeMock "microservice/internal/adapter/locale/mocks"
	loggerMock "microservice/internal/adapter/logger/mocks"
	"microservice/internal/core/domain"
	todoRepoMock "microservice/internal/core/port/mocks"
	"sync"
	"testing"
	"time"
)

func TestTodoUsecase_Create(t *testing.T) {
	id := uuid.New()
	description := "create mock item"
	datetime, _ := time.Parse(time.DateTime, "2025-08-07 10:11:12")

	testTodo := domain.NewTodo()
	testTodo.SetUUID(&id)
	testTodo.SetDescription(&description)
	testTodo.SetDueDate(&datetime)

	expectedTodo := domain.NewTodo()
	expectedTodo.SetUUID(&id)
	expectedTodo.SetDescription(&description)
	expectedTodo.SetDueDate(&datetime)

	t.Run("successful create with wait-group", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		//

		locale := localeMock.NewMockILocale(ctrl)
		logger := loggerMock.NewMockILogger(ctrl)
		todoRepo := todoRepoMock.NewMockITodoRepository(ctrl)

		//

		uc := NewTodo(logger, locale, todoRepo)

		//

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		wg := sync.WaitGroup{}
		wg.Add(1)

		todoRepo.EXPECT().Create(ctx, testTodo).Return(expectedTodo, nil).Times(1)

		result, err := uc.Create(ctx, testTodo)

		assert.NoError(t, err)
		assert.Equal(t, expectedTodo.UUID(), result.UUID())

		// wait for goroutine used in `Create` method to complete with timeout
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			t.Log("goroutine completed")
		case <-time.After(3 * time.Second):
			t.Fatal("goroutine timeout")
		}
	})

	t.Run("successful create with sleep for goroutine", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		//

		locale := localeMock.NewMockILocale(ctrl)
		logger := loggerMock.NewMockILogger(ctrl)
		todoRepo := todoRepoMock.NewMockITodoRepository(ctrl)

		//

		uc := NewTodo(logger, locale, todoRepo)

		//

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		todoRepo.EXPECT().Create(ctx, testTodo).Return(expectedTodo, nil).Times(1)

		result, err := uc.Create(ctx, testTodo)

		assert.NoError(t, err)
		assert.Equal(t, expectedTodo.UUID(), result.UUID())

		// sleep to allow goroutine to be completed
		time.Sleep(100 * time.Millisecond)

	})

	t.Run("todo repo error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		//

		locale := localeMock.NewMockILocale(ctrl)
		logger := loggerMock.NewMockILogger(ctrl)
		todoRepo := todoRepoMock.NewMockITodoRepository(ctrl)

		//

		uc := NewTodo(logger, locale, todoRepo)

		//

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		expectedErr := fmt.Errorf("repository error")

		todoRepo.EXPECT().Create(ctx, testTodo).Return(nil, expectedErr).Times(1)
		// no queue or logger expectations (goroutine won't run)

		result, err := uc.Create(ctx, testTodo)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, result)
	})

	t.Run("todo queue error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		//

		locale := localeMock.NewMockILocale(ctrl)
		logger := loggerMock.NewMockILogger(ctrl)
		todoRepo := todoRepoMock.NewMockITodoRepository(ctrl)

		//

		uc := NewTodo(logger, locale, todoRepo)

		//

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		wg := sync.WaitGroup{}
		wg.Add(1)

		queueErr := errors.New("queue consume failure")

		todoRepo.EXPECT().Create(ctx, testTodo).Return(expectedTodo, nil).Times(1)

		logger.EXPECT().Error("todo.uc.create.queue.send", gomock.Any()).DoAndReturn(func(scope string, fields ...interface{}) {
			errorFound := false
			errorItemIdFound := false

			for _, field := range fields {
				if zapField, ok := field.(zap.Field); ok {
					switch zapField.Key {
					case "error":
						if zapField.Interface != nil {
							if err, ok := zapField.Interface.(error); ok {
								assert.Equal(t, queueErr.Error(), err.Error(), "expected queue error")
								errorFound = true
							}
						}
					case "item.id":
						if zapField.String != "" {
							assert.Equal(t, expectedTodo.UUID().String(), zapField.String, "expected specific queue item id")
							errorItemIdFound = true
						}
					}
				}
			}

			assert.True(t, errorFound, "expected error found in log")
			assert.True(t, errorItemIdFound, "expected item id found in log")
		}).Times(1)

		result, err := uc.Create(ctx, testTodo)

		assert.NoError(t, err)
		assert.Equal(t, nil, err)
		assert.Equal(t, expectedTodo.UUID(), result.UUID())

		// wait for goroutine used in `Create` method to complete with timeout
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			t.Log("goroutine completed")
		case <-time.After(3 * time.Second):
			t.Fatal("goroutine timeout")
		}
	})
}
