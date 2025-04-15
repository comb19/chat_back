package usecase

import (
	"todo_back/domain/repository"

	"gorm.io/gorm"
)

type TodoUsecase interface {
	Insert(DB *gorm.DB, title, description string) error
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
