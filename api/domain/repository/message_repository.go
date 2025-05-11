package repository

import (
	"chat_back/domain/model"
)

type MessageRepository interface {
	Insert(channelID string, userID string, content string) (*model.Message, error)
	GetByID(ID string) (*model.Message, error)
	GetAllInChannel(channelID string) ([]model.Message, error)
}
