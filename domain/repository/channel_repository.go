package repository

import (
	"todo_back/domain/model"

	"gorm.io/gorm"
)

type ChannelRepository interface {
	Insert(DB *gorm.DB, name string, description string) error
	GetByID(DB *gorm.DB, id string) (model.Channel, error)
	GetAllInGuild(DB *gorm.DB, guild_id string) ([]model.Channel, error)
}
