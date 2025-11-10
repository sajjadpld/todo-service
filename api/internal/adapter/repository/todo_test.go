package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	localeMock "microservice/internal/adapter/locale/mocks"
	loggerMock "microservice/internal/adapter/logger/mocks"
	ormMock "microservice/internal/adapter/orm/mocks"
	"microservice/internal/adapter/orm/model"
	"microservice/internal/core/domain"
	"testing"
	"time"
)

func TestTodoRepository_Create(t *testing.T) {
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

	t.Run("successful create", func(t *testing.T) {
		dbConn, dbErr := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Silent),
		})

		if dbErr != nil {
			t.Fatalf("failed to open in-memory db: %v", dbErr)
		}

		if dbErr = dbConn.AutoMigrate(&model.Todos{}); dbErr != nil {
			t.Fatalf("failed to auto-migrate: %v", dbErr)
		}

		t.Cleanup(func() {
			sql, err := dbConn.DB()
			if err != nil {
				t.Logf("cleanup error: %v", err)
				return
			}

			if err = sql.Close(); err != nil {
				t.Log("sql conn close failure: ", err)
			}
		})

		//

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		locale := localeMock.NewMockILocale(ctrl)
		logger := loggerMock.NewMockILogger(ctrl)
		db := ormMock.NewMockISql(ctrl)

		ctx := context.Background()

		db.EXPECT().C().Return(dbConn)

		repo := NewTodo(locale, logger, db)
		res, err := repo.Create(ctx, testTodo)

		assert.Nil(t, err)
		assert.Equal(t, expectedTodo.UUID(), res.UUID())
		assert.GreaterOrEqual(t, time.Now(), res.CreatedAt())
	})

	t.Run("create failure", func(t *testing.T) {
		type TodoWithNewColumn struct {
			model.Todos
			Consume string `gorm:"not null"`
		}

		dbConn, dbErr := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Silent),
		})

		if dbErr != nil {
			t.Fatalf("failed to open in-memory db: %v", dbErr)
		}

		if dbErr = dbConn.AutoMigrate(&TodoWithNewColumn{}); dbErr != nil {
			t.Fatalf("failed to auto-migrate: %v", dbErr)
		}

		t.Cleanup(func() {
			sql, err := dbConn.DB()
			if err != nil {
				t.Logf("cleanup error: %v", err)
				return
			}

			if err = sql.Close(); err != nil {
				t.Log("sql conn close failure: ", err)
			}
		})

		//

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		locale := localeMock.NewMockILocale(ctrl)
		logger := loggerMock.NewMockILogger(ctrl)
		db := ormMock.NewMockISql(ctrl)

		ctx := context.Background()

		db.EXPECT().C().Return(dbConn)
		logger.EXPECT().Error("todo.repo.create", gomock.Any()).Times(1)

		repo := NewTodo(locale, logger, db)
		res, err := repo.Create(ctx, testTodo)
		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("create failure duplication error", func(t *testing.T) {
		type Todos struct {
			model.BaseSql
			Description string    `json:"description" gorm:"unique"`
			DueDate     time.Time `json:"dueDate"`
		}

		dbConn, dbErr := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Silent),
		})

		if dbErr != nil {
			t.Fatalf("failed to open in-memory db: %v", dbErr)
		}

		if dbErr = dbConn.AutoMigrate(&Todos{}); dbErr != nil {
			t.Fatalf("failed to auto-migrate: %v", dbErr)
		}

		t.Cleanup(func() {
			sql, err := dbConn.DB()
			if err != nil {
				t.Logf("cleanup error: %v", err)
				return
			}

			if err = sql.Close(); err != nil {
				t.Log("sql conn close failure: ", err)
			}
		})

		//

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		locale := localeMock.NewMockILocale(ctrl)
		logger := loggerMock.NewMockILogger(ctrl)
		db := ormMock.NewMockISql(ctrl)

		ctx := context.Background()

		db.EXPECT().C().Return(dbConn).Times(2)
		logger.EXPECT().Error("todo.repo.create", gomock.Any()).Times(1)

		repo := NewTodo(locale, logger, db)

		// first insert: should succeed
		_, err := repo.Create(ctx, testTodo)
		if err != nil {
			t.Fatalf("duplication err test - first insert err %s", err)
		}

		// second insert: duplicate error
		_, err = repo.Create(ctx, testTodo)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "UNIQUE constraint failed: todos.description")

	})
}
