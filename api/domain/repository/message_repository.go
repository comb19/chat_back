package repository

import (
	"chat_back/domain/model"

	"gorm.io/gorm"
)

type MessageRepository interface {
	Insert(db *gorm.DB, channelID string, userID string, content string) error
	GetByID(db *gorm.DB, ID string) (model.Message, error)
	GetAllInChannel(db *gorm.DB, channelID string) ([]model.Message, error)
}
