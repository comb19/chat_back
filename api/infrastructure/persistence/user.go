package persistence

import (
	"chat_back/domain/model"
	"chat_back/domain/repository"

	"gorm.io/gorm"
)

type userPersistence struct {
	db *gorm.DB
}

func NewUserPersistence(db *gorm.DB) repository.UserRepository {
	return &userPersistence{
		db: db,
	}
}

func (up userPersistence) Insert(id, userName string) (*model.User, error) {
	user := &model.User{
		ID:       id,
		UserName: userName,
	}

	if err := up.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (up userPersistence) Find(id string) (*model.User, error) {
	user := &model.User{}
	if err := up.db.First(user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return user, nil
}
