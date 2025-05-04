package persistence

import (
	"chat_back/domain/model"
	"chat_back/domain/repository"
	"fmt"

	"gorm.io/gorm"
)

type channelPersistence struct{}

func NewChannelPersistence() repository.ChannelRepository {
	return &channelPersistence{}
}

func (c *channelPersistence) Insert(db *gorm.DB, name string, description string, private bool, guildID *string) (*string, error) {
	channel := model.Channel{
		Name:        name,
		Description: description,
		Private:     private,
		GuildID:     guildID,
	}
	result := db.Select("name", "description", "private", "guild_id").Create(&channel)
	fmt.Println(result)
	fmt.Println(channel)
	if result.Error != nil {
		return nil, result.Error
	}
	return &channel.ID, result.Error
}

func (c *channelPersistence) GetByID(db *gorm.DB, id string) (model.Channel, error) {
	var channel model.Channel
	result := db.First(&channel, id)
	if result.Error != nil {
		return model.Channel{}, result.Error
	}
	return channel, nil
}

func (c *channelPersistence) GetAllInGuild(db *gorm.DB, guildID *string) ([]model.Channel, error) {
	var channels []model.Channel
	result := db.Where("guild_id = ?", *guildID).Find(&channels)
	if result.Error != nil {
		return nil, result.Error
	}
	return channels, nil
}
