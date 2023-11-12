package repository

import (
	"context"
	"database/sql"
	"github.com/Ikhlashmulya/golang-grpc-server/entity"
	"github.com/Ikhlashmulya/golang-grpc-server/exception"
)

type TodoRepositoryImpl struct {
	DB *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepositoryImpl {
	return &TodoRepositoryImpl{
		DB: db,
	}
}

func (todorepositoryimpl *TodoRepositoryImpl) Create(ctx context.Context, todo entity.Todo) {
	tx, err := todorepositoryimpl.DB.Begin()
	exception.PanicIfError(err)
	defer commitOrRollback(tx)

	SQL := "INSERT INTO todo (id, name) VALUES (?, ?)"

	_, err = tx.ExecContext(ctx, SQL, todo.Id, todo.Name)
	exception.PanicIfError(err)
}

func (todorepositoryimpl *TodoRepositoryImpl) GetAll(ctx context.Context) []entity.Todo {
	SQL := "SELECT id, name FROM todo"

	rows, err := todorepositoryimpl.DB.QueryContext(ctx, SQL)
	exception.PanicIfError(err)
	defer rows.Close()

	todos := []entity.Todo{}

	for rows.Next() {
		todo := entity.Todo{}

		errScan := rows.Scan(&todo.Id, &todo.Name)
		exception.PanicIfError(errScan)

		todos = append(todos, todo)
	}

	return todos
}

func (todorepositoryimpl *TodoRepositoryImpl) Delete(ctx context.Context, todoId string) {
	tx, err := todorepositoryimpl.DB.Begin()
	exception.PanicIfError(err)
	defer commitOrRollback(tx)

	SQL := "DELETE FROM todo WHERE id = ?"

	_, err = tx.ExecContext(ctx, SQL, todoId)
	exception.PanicIfError(err)
}

func (todorepositoryimpl *TodoRepositoryImpl) GetById(ctx context.Context, todoId string) (response entity.Todo, err error) {
	SQL := "SELECT id, name FROM todo WHERE id = ?"

	row, err := todorepositoryimpl.DB.QueryContext(ctx, SQL, todoId)
	exception.PanicIfError(err)
	defer row.Close()

	if row.Next() {
		errScan := row.Scan(&response.Id, &response.Name)
		exception.PanicIfError(errScan)

		return response, nil
	} else {
		return entity.Todo{}, sql.ErrNoRows
	}
}

func commitOrRollback(tx *sql.Tx) {
	err := recover()
	switch err {
	case nil:
		errCommit := tx.Commit()
		exception.PanicIfError(errCommit)
	default:
		errRollback := tx.Rollback()
		exception.PanicIfError(errRollback)
		panic(errRollback)
	}
}
