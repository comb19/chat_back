package persistence

import (
	"chat_back/domain/model"
	"chat_back/domain/repository"

	"gorm.io/gorm"
)

type messagePersistence struct{}

func NewMessagePersistence() repository.MessageRepository {
	return &messagePersistence{}
}

func (mp messagePersistence) Insert(db *gorm.DB, channelID string, userID string, content string) error {
	message := model.Message{
		Channel_id: channelID,
		User_id:    userID,
		Content:    content,
	}
	result := db.Create(&message)
	return result.Error
}

func (mp messagePersistence) GetByID(db *gorm.DB, ID string) (model.Message, error) {
	var message model.Message
	result := db.First(&message, ID)
	if result.Error != nil {
		return model.Message{}, result.Error
	}
	return message, nil
}
func (mp messagePersistence) GetAllInChannel(db *gorm.DB, channelID string) ([]model.Message, error) {
	var messages []model.Message
	result := db.Where("channel_id = ?", channelID).Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}
