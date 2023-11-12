package repository

import (
	"context"
	"github.com/Ikhlashmulya/golang-grpc-server/entity"
)

type TodoRepository interface {
	Create(ctx context.Context, todo entity.Todo)
	GetAll(ctx context.Context) []entity.Todo
	GetById(ctx context.Context, todoId string) (response entity.Todo, err error)
	Delete(ctx context.Context, todoId string)
}
