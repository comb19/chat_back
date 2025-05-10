package persistence

import (
	"chat_back/domain/model"
	"chat_back/domain/repository"
	"fmt"

	"gorm.io/gorm"
)

type userChannelsPersistence struct {
	db *gorm.DB
}

func NewUserChannelsPersistence(db *gorm.DB) repository.UserChannelsRepository {
	return &userChannelsPersistence{
		db: db,
	}
}

func (c *userChannelsPersistence) Insert(db *gorm.DB, userID string, channelID string) (*model.UserChannels, error) {
	fmt.Println("Inserting user channel:", userID, channelID)
	if userChannels, err := c.Find(db, userID, channelID); err != nil {
		return nil, err
	} else if userChannels != nil {
		return userChannels, fmt.Errorf("user %s is already a member of channel %s", userID, channelID)
	}
	userChannels := &model.UserChannels{
		UserID:    userID,
		ChannelID: channelID,
	}
	if err := db.Create(userChannels).Error; err != nil {
		return nil, err
	}
	fmt.Println("User channel created:", userChannels)
	return userChannels, nil
}

func (c *userChannelsPersistence) Find(db *gorm.DB, userID string, channelID string) (*model.UserChannels, error) {
	var userChannels model.UserChannels
	result := db.Where("user_id = ? AND channel_id = ?", userID, channelID).FirstOrInit(&userChannels)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		fmt.Println("User channel not found")
		return nil, nil
	}
	return &userChannels, nil
}
