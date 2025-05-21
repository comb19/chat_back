package repository

import (
	"chat_back/domain/model"
)

type UserRepository interface {
	Insert(id, userName string) (*model.User, error)
	Find(id string) (*model.User, error)
}
