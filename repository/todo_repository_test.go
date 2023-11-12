package repository

import (
	"context"
	"database/sql"
	"github.com/Ikhlashmulya/golang-grpc-server/entity"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.Nil(t, err)

		todoRepository := NewTodoRepository(db)

		todo := entity.Todo{
			Id:   "1",
			Name: "test",
		}

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO todo").WithArgs(todo.Id, todo.Name).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		todoRepository.Create(context.Background(), todo)

		assert.Nil(t, mock.ExpectationsWereMet())

	})

	t.Run("Fail", func(t *testing.T) {
		assert.Panics(t, func() {
			db, mock, err := sqlmock.New()
			assert.Nil(t, err)

			todoRepository := NewTodoRepository(db)

			todo := entity.Todo{
				Id:   "1",
				Name: "test",
			}

			mock.ExpectBegin()
			mock.ExpectExec("INSERT INTO todo").WithArgs(todo.Id, todo.Name).WillReturnError(sqlmock.ErrCancelled)
			mock.ExpectRollback()

			todoRepository.Create(context.Background(), todo)

			assert.Nil(t, mock.ExpectationsWereMet())

		})
	})
}

func TestDelete(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.Nil(t, err)

		todoRepository := NewTodoRepository(db)

		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM todo").WithArgs("1").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		todoRepository.Delete(context.Background(), "1")

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		assert.Panics(t, func() {
			db, mock, err := sqlmock.New()
			assert.Nil(t, err)

			todoRepository := NewTodoRepository(db)

			mock.ExpectBegin()
			mock.ExpectExec("DELETE FROM todo").WithArgs("1").WillReturnError(sql.ErrNoRows)
			mock.ExpectRollback()

			todoRepository.Delete(context.Background(), "1")

			assert.Nil(t, mock.ExpectationsWereMet())

		})
	})

}

func TestGetById(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.Nil(t, err)

		todoRepository := NewTodoRepository(db)

		todo := entity.Todo{
			Id:   "1",
			Name: "test",
		}

		mock.ExpectQuery("SELECT").WithArgs("1").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(todo.Id, todo.Name))

		gotTodo, err := todoRepository.GetById(context.Background(), "1")
		assert.Nil(t, err)
		assert.Equal(t, todo.Id, gotTodo.Id)
		assert.Equal(t, todo.Name, gotTodo.Name)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.Nil(t, err)

		todoRepository := NewTodoRepository(db)

		mock.ExpectQuery("SELECT").WithArgs("2").WillReturnRows(&sqlmock.Rows{})

		goTodo, err := todoRepository.GetById(context.Background(), "2")
		if assert.NotNil(t, err) {
			assert.Equal(t, sql.ErrNoRows, err)
		}
		assert.Empty(t, goTodo)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	todoRepository := NewTodoRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow("1", "test1").AddRow("2", "test2")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	gotTodos := todoRepository.GetAll(context.Background())

	assert.NotEmpty(t, gotTodos)

	assert.Nil(t, mock.ExpectationsWereMet())
}
