package usecase

import (
	"chat_back/domain/model"
	"chat_back/domain/repository"

	"gorm.io/gorm"
)

type userUseCase struct {
	userRepository repository.UserRepository
}

type UserUsecase interface {
	CreateUserByClerk(db *gorm.DB, id, userName string) (*model.User, error)
	GetUserByID(db *gorm.DB, id string) (*model.User, error)
}

func NewUserUsecase(userRepository repository.UserRepository) UserUsecase {
	return &userUseCase{
		userRepository: userRepository,
	}
}

func (uu *userUseCase) CreateUserByClerk(db *gorm.DB, id, userName string) (*model.User, error) {
	user, err := uu.userRepository.Insert(db, id, userName)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uu *userUseCase) GetUserByID(db *gorm.DB, id string) (*model.User, error) {
	user, err := uu.userRepository.GetByID(db, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
