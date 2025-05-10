package repository

import (
	"chat_back/domain/model"
)

type UserChannelsRepository interface {
	Insert(userID string, channelID string) (*model.UserChannels, error)
	Find(userID string, channelID string) (*model.UserChannels, error)
}
