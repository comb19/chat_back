package repository

import "gorm.io/gorm"

type TodoRepository interface {
	Insert(DB *gorm.DB, title, description string) error
	// GetAll() ([]model.Todo, error)
	// GetByID(id int) (model.Todo, error)
	// Update(id int, title, description string) error
	// Delete(id int) error
}
