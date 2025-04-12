package repository

import (
	"todo_back/domain/model"
)

type TodoRepository interface {
	Insert(title, description string) error
	GetAll() ([]model.Todo, error)
	GetByID(id int) (model.Todo, error)
	Update(id int, title, description string) error
	Delete(id int) error
}
