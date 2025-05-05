package repository

import (
	"chat_back/domain/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Insert(DB *gorm.DB, id, userName string) (*model.User, error)
	GetByID(DB *gorm.DB, id string) (*model.User, error)
}
