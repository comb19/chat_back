package repository

import (
	"chat_back/domain/model"
)

type UserRepository interface {
	Insert(id, userName string) (*model.User, error)
	GetByID(id string) (*model.User, error)
}
