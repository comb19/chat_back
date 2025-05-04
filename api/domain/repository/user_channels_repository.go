package repository

import (
	"chat_back/domain/model"

	"gorm.io/gorm"
)

type UserChannelsRepository interface {
	Insert(db *gorm.DB, userID string, channelID string) (*model.UserChannels, error)
	Find(db *gorm.DB, userID string, channelID string) (*model.UserChannels, error)
}
