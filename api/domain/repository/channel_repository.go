package repository

import (
	"chat_back/domain/model"
)

type ChannelRepository interface {
	Insert(name string, description string, private bool, guildID *string) (*model.Channel, error)
	Delete(id string) error
	Find(id string) (*model.Channel, error)
	FindAllInGuild(guildID *string) ([]*model.Channel, error)
}
