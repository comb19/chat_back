package persistence

import (
	"chat_back/domain/model"
	"chat_back/domain/repository"
	"errors"
	"fmt"
	"log/slog"

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

func (ccp *userChannelsPersistence) Insert(userID string, channelID string) (*model.UserChannels, error) {
	slog.Debug(fmt.Sprintf("Inserting user: %s, channel: %s", userID, channelID))

	if userChannels, err := ccp.Find(userID, channelID); err != nil {
		return nil, err
	} else if userChannels != nil {
		return userChannels, fmt.Errorf("user %s is already a member of channel %s", userID, channelID)
	}

	userChannels := &model.UserChannels{
		UserID:    userID,
		ChannelID: channelID,
	}
	if err := ccp.db.Create(userChannels).Error; err != nil {
		return nil, err
	}
	slog.Debug(fmt.Sprint("User channel created:", userChannels))

	return userChannels, nil
}

func (ccp *userChannelsPersistence) Find(userID string, channelID string) (*model.UserChannels, error) {
	var userChannels model.UserChannels
	result := ccp.db.Where("user_id = ? AND channel_id = ?", userID, channelID).First(&userChannels)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Info(fmt.Sprintf("a relation between user %s and channe %s is not found.", userID, channelID))
			return nil, nil
		}
		slog.Error(err.Error())
		return nil, err
	}
	return &userChannels, nil
}
