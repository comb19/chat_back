package repository

import "chat_back/domain/model"

type GuildRepository interface {
	Insert(name, description string) (*model.Guild, error)
	Find(id string) (*model.Guild, error)
	FindOfUser(userId string) ([]*model.Guild, error)
}
