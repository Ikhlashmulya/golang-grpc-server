package grpc

import (
	"context"
	"testing"

	"github.com/Ikhlashmulya/golang-grpc-server/model"
	pb "github.com/Ikhlashmulya/golang-grpc-server/delivery/grpc/proto"
	"github.com/Ikhlashmulya/golang-grpc-server/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	todoUsecase := usecase.NewMockTodoUsecase(ctrl)

	todoUsecase.EXPECT().Create(gomock.Any(), model.CreateTodoRequest{Name: "Testing"}).
		Return(model.TodoResponse{Id: "1", Name: "Testing"})

	todoHandler := NewTodoHandler(todoUsecase)

	response, err := todoHandler.Create(context.Background(), &pb.CreateTodoRequest{Name: "Testing"})
	assert.Nil(t, err)

	assert.Equal(t, "1", response.GetTodo().GetId())
	assert.Equal(t, "Testing", response.GetTodo().GetName())
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	todoUsecase := usecase.NewMockTodoUsecase(ctrl)

	todoUsecase.EXPECT().Delete(gomock.Any(), "1").Times(1)

	todoHandler := NewTodoHandler(todoUsecase)

	_, err := todoHandler.Delete(context.Background(), &pb.DeleteTodoRequest{Id: "1"})
	assert.Nil(t, err)
	assert.True(t, ctrl.Satisfied())
}

func TestGetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	todoUsecase := usecase.NewMockTodoUsecase(ctrl)

	todoUsecase.EXPECT().GetById(gomock.Any(), "1").
		Return(model.TodoResponse{Id: "1", Name: "Testing"})

	todoHandler := NewTodoHandler(todoUsecase)

	response, err := todoHandler.GetById(context.Background(), &pb.GetByIdTodoRequest{Id: "1"})
	assert.Nil(t, err)

	assert.Equal(t, "1", response.GetTodo().GetId())
	assert.Equal(t, "Testing", response.GetTodo().GetName())
}

func TestGetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	todoUsecase := usecase.NewMockTodoUsecase(ctrl)

	todoResponse := []model.TodoResponse{
		{
			Id:   "1",
			Name: "Testing 1",
		},
		{
			Id:   "2",
			Name: "Testing 2",
		},
	}

	todoUsecase.EXPECT().GetAll(gomock.Any()).
		Return(todoResponse)

	todoHandler := NewTodoHandler(todoUsecase)

	responses, err := todoHandler.GetAll(context.Background(), &pb.GetAllTodoRequest{})
	assert.Nil(t, err)

	for i := 0; i < len(responses.GetTodo()); i++ {
		assert.Equal(t, todoResponse[i].Id, responses.GetTodo()[i].GetId())
		assert.Equal(t, todoResponse[i].Name, responses.GetTodo()[i].GetName())
	}
}
