package grpc

import (
	"context"

	"github.com/Ikhlashmulya/golang-grpc-server/model"
	pb "github.com/Ikhlashmulya/golang-grpc-server/delivery/grpc/proto"
	"github.com/Ikhlashmulya/golang-grpc-server/usecase"
)

type TodoHandler struct {
	TodoUsecase usecase.TodoUsecase
	pb.UnimplementedTodoServiceServer
}

func NewTodoHandler(todoUsecase usecase.TodoUsecase) *TodoHandler {
	return &TodoHandler{TodoUsecase: todoUsecase}
}

func (todohandler *TodoHandler) Create(ctx context.Context, request *pb.CreateTodoRequest) (*pb.CreateTodoResponse, error) {
	result := todohandler.TodoUsecase.Create(ctx, model.CreateTodoRequest{
		Name: request.GetName(),
	})

	return &pb.CreateTodoResponse{
		Todo: modelToProto(result),
	}, nil
}

func (todohandler *TodoHandler) Delete(ctx context.Context, request *pb.DeleteTodoRequest) (*pb.DeleteTodoResponse, error) {
	todohandler.TodoUsecase.Delete(ctx, request.GetId())

	return &pb.DeleteTodoResponse{}, nil
}

func (todohandler *TodoHandler) GetAll(ctx context.Context, request *pb.GetAllTodoRequest) (*pb.GetAllTodoResponse, error) {
	results := todohandler.TodoUsecase.GetAll(ctx)

	var todo []*pb.Todo
	for _, result := range results {
		todo = append(todo, modelToProto(result))
	}

	return &pb.GetAllTodoResponse{Todo: todo}, nil
}

func (todohandler *TodoHandler) GetById(ctx context.Context, request *pb.GetByIdTodoRequest) (*pb.GetByIdTodoResponse, error) {
	result := todohandler.TodoUsecase.GetById(ctx, request.GetId())

	return &pb.GetByIdTodoResponse{Todo: modelToProto(result)}, nil
}

func modelToProto(model model.TodoResponse) (*pb.Todo) {
	return &pb.Todo{
		Id:   model.Id,
		Name: model.Name,
	}
}
