package usecase

import (
	"chat_back/domain/model"
	"chat_back/domain/repository"
)

type userUseCase struct {
	userRepository repository.UserRepository
}

type UserUsecase interface {
	CreateUserByClerk(id, userName string) (*model.User, error)
	GetUserByID(id string) (*model.User, error)
}

func NewUserUsecase(userRepository repository.UserRepository) UserUsecase {
	return &userUseCase{
		userRepository: userRepository,
	}
}

func (uu *userUseCase) CreateUserByClerk(id, userName string) (*model.User, error) {
	user, err := uu.userRepository.Insert(id, userName)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uu *userUseCase) GetUserByID(id string) (*model.User, error) {
	user, err := uu.userRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
