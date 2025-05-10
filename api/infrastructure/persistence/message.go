package persistence

import (
	"chat_back/domain/model"
	"chat_back/domain/repository"
	"fmt"

	"gorm.io/gorm"
)

type messagePersistence struct {
	db *gorm.DB
}

func NewMessagePersistence(db *gorm.DB) repository.MessageRepository {
	return &messagePersistence{
		db: db,
	}
}

func (mp messagePersistence) Insert(channelID string, userID string, content string) (*model.Message, error) {
	message := model.Message{
		ChannelID: channelID,
		UserID:    userID,
		Content:   content,
	}
	result := mp.db.Select("user_id", "channel_id", "content").Create(&message)
	return &message, result.Error
}

func (mp messagePersistence) GetByID(ID string) (*model.Message, error) {
	var message model.Message
	result := mp.db.First(&message, ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &message, nil
}
func (mp messagePersistence) GetAllInChannel(channelID string) ([]model.Message, error) {
	fmt.Println("GetAllInChannel")
	var messages []model.Message
	result := mp.db.Table("messages").Select("messages.id, messages.user_id, users.user_name, messages.content, messages.channel_id").Where("channel_id = ?", channelID).Joins("left outer join users on messages.user_id = users.id").Scan(&messages)
	if result.Error != nil {
		return nil, result.Error
	}
	fmt.Println(messages)
	return messages, nil
}
