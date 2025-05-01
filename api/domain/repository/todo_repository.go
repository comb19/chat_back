package repository

import (
	"chat_back/domain/model"

	"gorm.io/gorm"
)

type TodoRepository interface {
	Insert(DB *gorm.DB, title, description string) error
	GetAll(DB *gorm.DB) ([]model.Todo, error)
	// GetByID(id int) (model.Todo, error)
	// Update(id int, title, description string) error
	// Delete(id int) error
}
