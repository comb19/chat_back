package repository

import (
	"chat_back/domain/model"

	"gorm.io/gorm"
)

type ChannelRepository interface {
	Insert(db *gorm.DB, name string, description string, private bool, guildID *string) (*string, error)
	GetByID(db *gorm.DB, id string) (model.Channel, error)
	GetAllInGuild(db *gorm.DB, guildID *string) ([]model.Channel, error)
}
