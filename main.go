package main

import (
	"net"

	"github.com/Ikhlashmulya/golang-grpc-server/app"
	handler "github.com/Ikhlashmulya/golang-grpc-server/delivery/grpc"
	"github.com/Ikhlashmulya/golang-grpc-server/exception"
	pb "github.com/Ikhlashmulya/golang-grpc-server/delivery/grpc/proto"
	"github.com/Ikhlashmulya/golang-grpc-server/repository"
	"github.com/Ikhlashmulya/golang-grpc-server/usecase"
	"google.golang.org/grpc"
)

func main() {
	db := app.NewDB()
	todoRepository := repository.NewTodoRepository(db)
	todoUsecase := usecase.NewTodoUsecase(todoRepository)
	todoHandler := handler.NewTodoHandler(todoUsecase)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(exception.PanicHandlerInterceptor()),
	)

	pb.RegisterTodoServiceServer(server, todoHandler)

	listener, err := net.Listen("tcp", ":50051")
	exception.PanicIfError(err)

	err = server.Serve(listener)
	exception.PanicIfError(err)
}
