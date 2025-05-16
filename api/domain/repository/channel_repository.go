package repository

import (
	"chat_back/domain/model"
)

type ChannelRepository interface {
	Insert(name string, description string, private bool, guildID *string) (*model.Channel, error)
	GetByID(id string) (*model.Channel, error)
	GetAllInGuild(guildID *string) ([]*model.Channel, error)
}
