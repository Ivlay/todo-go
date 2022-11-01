package repository

import (
	"github.com/Ivlay/todo-go"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GetUser(username, password string) (todo.User, error)
}

type TodoList interface {
	Create(userId int, list todo.TodoList) (int, error)
	GetAll(userId int) ([]todo.TodoList, error)
	GetById(userId, listId int) (todo.TodoList, error)
	Update(userId, listId int, input todo.UpdateListInput) error
	Delete(userId, listId int) error
}

type TodoItem interface {
	Create(listId int, item todo.TodoItem) (int, error)
	GetAll(userId, listId int) ([]todo.TodoItem, error)
	GetById(userId, itemId int) (todo.TodoItem, error)
	Update(userId, itemId int, input todo.UpdateItemInput) error
	Delete(userId, itemId int) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRspository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthSql(db),
		TodoList: NewTodoListSql(db),
		TodoItem: NewTodoItemSql(db),
	}
}