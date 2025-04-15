package persistence

import (
	"todo_back/domain/model"
	"todo_back/domain/repository"

	"gorm.io/gorm"
)

type todoPersistence struct{}

func NewTodoPersistence() repository.TodoRepository {
	return &todoPersistence{}
}

func (tp todoPersistence) Insert(DB *gorm.DB, title, description string) error {
	todo := model.Todo{
		Title:       title,
		Description: description}
	result := DB.Create(&todo)
	return result.Error
}
