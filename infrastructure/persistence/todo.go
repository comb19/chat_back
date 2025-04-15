package persistence

import (
	"fmt"
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
	fmt.Println("todo", todo)
	result := DB.Create(&todo)
	return result.Error
}

func (tp todoPersistence) GetAll(DB *gorm.DB) ([]model.Todo, error) {
	var todos []model.Todo
	result := DB.Find(&todos)
	if result.Error != nil {
		return nil, result.Error
	}
	return todos, nil
}
