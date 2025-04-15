package usecase

import (
	"todo_back/domain/model"
	"todo_back/domain/repository"

	"gorm.io/gorm"
)

type TodoUsecase interface {
	Insert(DB *gorm.DB, title, description string) error
	GetAll(DB *gorm.DB) ([]model.Todo, error)
}

type todoUseCase struct {
	todoRepository repository.TodoRepository
}

func NewTodoUsecase(todoRepository repository.TodoRepository) TodoUsecase {
	return &todoUseCase{
		todoRepository: todoRepository,
	}
}

func (tu todoUseCase) Insert(DB *gorm.DB, title, description string) error {
	err := tu.todoRepository.Insert(DB, title, description)
	if err != nil {
		return err
	}
	return nil
}

func (tu todoUseCase) GetAll(DB *gorm.DB) ([]model.Todo, error) {
	todos, err := tu.todoRepository.GetAll(DB)
	if err != nil {
		return nil, err
	}
	return todos, nil
}
