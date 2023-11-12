package test

import (
	"context"
	"database/sql"
	"net"
	"testing"

	"github.com/Ikhlashmulya/golang-grpc-server/app"
	handler "github.com/Ikhlashmulya/golang-grpc-server/delivery/grpc"
	"github.com/Ikhlashmulya/golang-grpc-server/entity"
	"github.com/Ikhlashmulya/golang-grpc-server/exception"
	pb "github.com/Ikhlashmulya/golang-grpc-server/delivery/grpc/proto"
	"github.com/Ikhlashmulya/golang-grpc-server/repository"
	"github.com/Ikhlashmulya/golang-grpc-server/usecase"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	db             = app.NewDB()
	todoRepository = repository.NewTodoRepository(db)
	todoUsecase    = usecase.NewTodoUsecase(todoRepository)
)

func TestCreateTodo(t *testing.T) {
	server := setupGRPCServer()
	defer server.Stop()

	conn, err := setupGRPCClient()
	assert.Nil(t, err)
	defer conn.Close()

	client := pb.NewTodoServiceClient(conn)

	response, err := client.Create(context.Background(), &pb.CreateTodoRequest{Name: "Do Testing"})
	assert.Nil(t, err)
	assert.NotEmpty(t, response)

	_, err = todoRepository.GetById(context.Background(), response.GetTodo().GetId())
	assert.Nil(t, err)

	todoRepository.Delete(context.Background(), response.GetTodo().GetId())
}

func TestDeleteTodo(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		todoRepository.Create(context.Background(), entity.Todo{
			Id:   "1",
			Name: "Do Testing",
		})

		server := setupGRPCServer()
		defer server.Stop()

		conn, err := setupGRPCClient()
		assert.Nil(t, err)
		defer conn.Close()

		client := pb.NewTodoServiceClient(conn)

		_, err = client.Delete(context.Background(), &pb.DeleteTodoRequest{Id: "1"})
		assert.Nil(t, err)

		_, err = todoRepository.GetById(context.Background(), "1")
		assert.NotNil(t, err)
	})

	t.Run("Not Found", func(t *testing.T) {
		server := setupGRPCServer()
		defer server.Stop()

		conn, err := setupGRPCClient()
		assert.Nil(t, err)
		defer conn.Close()

		client := pb.NewTodoServiceClient(conn)

		_, err = client.Delete(context.Background(), &pb.DeleteTodoRequest{Id: "1"})
		assert.NotNil(t, err)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, statusErr.Code())
		assert.Equal(t, sql.ErrNoRows.Error(), statusErr.Message())
	})
}

func TestGetByIdTodo(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		todoRepository.Create(context.Background(), entity.Todo{
			Id:   "1",
			Name: "Do Testing",
		})

		server := setupGRPCServer()
		defer server.Stop()

		conn, err := setupGRPCClient()
		assert.Nil(t, err)
		defer conn.Close()

		client := pb.NewTodoServiceClient(conn)

		response, err := client.GetById(context.Background(), &pb.GetByIdTodoRequest{Id: "1"})
		assert.Nil(t, err)
		assert.Equal(t, "1", response.GetTodo().GetId())
		assert.Equal(t, "Do Testing", response.GetTodo().GetName())

		todoRepository.Delete(context.Background(), "1")
	})

	t.Run("not found", func(t *testing.T) {
		server := setupGRPCServer()
		defer server.Stop()

		conn, err := setupGRPCClient()
		assert.Nil(t, err)
		defer conn.Close()

		client := pb.NewTodoServiceClient(conn)

		response, err := client.GetById(context.Background(), &pb.GetByIdTodoRequest{Id: "2"})
		assert.Nil(t, response)
		assert.NotNil(t, err)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, statusErr.Code())
		assert.Equal(t, sql.ErrNoRows.Error(), statusErr.Message())
	})
}

func TestGetAllTodo(t *testing.T) {
	todoRepository.Create(context.Background(), entity.Todo{
		Id:   "1",
		Name: "Do Testing",
	})

	todoRepository.Create(context.Background(), entity.Todo{
		Id:   "2",
		Name: "Do Coding",
	})

	server := setupGRPCServer()
	defer server.Stop()

	conn, err := setupGRPCClient()
	assert.Nil(t, err)
	defer conn.Close()

	client := pb.NewTodoServiceClient(conn)

	response, err := client.GetAll(context.Background(), &pb.GetAllTodoRequest{})
	assert.Nil(t, err)

	assert.Equal(t, 2, len(response.GetTodo()))

	todoRepository.Delete(context.Background(), "1")
	todoRepository.Delete(context.Background(), "2")
}

func setupGRPCServer() *grpc.Server {
	todoHandler := handler.NewTodoHandler(todoUsecase)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(exception.PanicHandlerInterceptor()),
	)

	pb.RegisterTodoServiceServer(server, todoHandler)

	listener, _ := net.Listen("tcp", ":50051")
	go func() {
		_ = server.Serve(listener)
	}()

	return server
}

func setupGRPCClient() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	return conn, err
}
