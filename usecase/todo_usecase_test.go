package usecase

import (
	"context"
	"database/sql"
	"github.com/Ikhlashmulya/golang-grpc-server/entity"
	"github.com/Ikhlashmulya/golang-grpc-server/model"
	"github.com/Ikhlashmulya/golang-grpc-server/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	todoRepository := repository.NewMockTodoRepository(ctrl)

	todoRepository.EXPECT().Create(context.Background(), gomock.Any()).Times(1)

	todoUsecase := NewTodoUsecase(todoRepository)
	response := todoUsecase.Create(context.Background(), model.CreateTodoRequest{Name: "test"})

	assert.NotEmpty(t, response)
	assert.True(t, ctrl.Satisfied())
}

func TestDelete(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		todoRepository := repository.NewMockTodoRepository(ctrl)

		todo := entity.Todo{
			Id:   "1",
			Name: "test",
		}

		todoRepository.EXPECT().GetById(context.Background(), todo.Id).Return(todo, nil)
		todoRepository.EXPECT().Delete(context.Background(), todo.Id).Times(1)

		todoUsecase := NewTodoUsecase(todoRepository)

		todoUsecase.Delete(context.Background(), "1")

		assert.True(t, ctrl.Satisfied())
	})
	t.Run("not found", func(t *testing.T) {
		assert.Panics(t, func() {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			todoRepository := repository.NewMockTodoRepository(ctrl)

			todo := entity.Todo{
				Id:   "1",
				Name: "test",
			}

			todoRepository.EXPECT().GetById(context.Background(), todo.Id).Return(todo, sql.ErrNoRows)
			// todoRepository.EXPECT().Delete(context.Background(), todo.Id).Times(1)

			todoUsecase := NewTodoUsecase(todoRepository)

			todoUsecase.Delete(context.Background(), "1")

			assert.True(t, ctrl.Satisfied())
		})
	})
}

func TestGetById(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		todoRepository := repository.NewMockTodoRepository(ctrl)

		todo := entity.Todo{
			Id:   "1",
			Name: "test",
		}

		todoRepository.EXPECT().GetById(context.Background(), todo.Id).Return(todo, nil)

		todoUsecase := NewTodoUsecase(todoRepository)

		response := todoUsecase.GetById(context.Background(), "1")

		assert.Equal(t, todo.Name, response.Name)
		assert.True(t, ctrl.Satisfied())
	})
	t.Run("not found", func(t *testing.T) {
		assert.Panics(t, func() {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			todoRepository := repository.NewMockTodoRepository(ctrl)

			todo := entity.Todo{
				Id:   "1",
				Name: "test",
			}

			todoRepository.EXPECT().GetById(context.Background(), todo.Id).Return(todo, sql.ErrNoRows)

			todoUsecase := NewTodoUsecase(todoRepository)

			_ = todoUsecase.GetById(context.Background(), "1")

			assert.True(t, ctrl.Satisfied())
		})
	})
}

func TestGetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	todoRepository := repository.NewMockTodoRepository(ctrl)

	todos := []entity.Todo{
		{
			Id:   "1",
			Name: "test1",
		},
		{
			Id:   "2",
			Name: "test2",
		},
	}

	todoRepository.EXPECT().GetAll(context.Background()).Return(todos)

	todoUsecase := NewTodoUsecase(todoRepository)

	responses := todoUsecase.GetAll(context.Background())

	assert.NotEmpty(t, responses)

	assert.True(t, ctrl.Satisfied())
}
