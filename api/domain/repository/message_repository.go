package repository

import (
	"chat_back/domain/model"
)

type MessageRepository interface {
	Insert(channelID string, userID string, content string) (*model.Message, error)
	Find(ID string) (*model.Message, error)
	FindAllInChannel(channelID string) (*[]*model.Message, error)
}
