package persistence

import (
	"chat_back/domain/model"
	"chat_back/domain/repository"
	"fmt"

	"gorm.io/gorm"
)

type channelPersistence struct {
	db *gorm.DB
}

func NewChannelPersistence(db *gorm.DB) repository.ChannelRepository {
	return &channelPersistence{
		db: db,
	}
}

func (cp *channelPersistence) Insert(name string, description string, private bool, guildID *string) (*model.Channel, error) {
	channel := model.Channel{
		Name:        name,
		Description: description,
		Private:     private,
		GuildID:     guildID,
	}
	result := cp.db.Select("name", "description", "private", "guild_id").Create(&channel)
	fmt.Println(result)
	fmt.Println(channel)
	if result.Error != nil {
		return nil, result.Error
	}
	fmt.Println("Channel created:", channel)
	return &channel, result.Error
}

func (cp *channelPersistence) GetByID(id string) (*model.Channel, error) {
	var channel model.Channel
	result := cp.db.Where("id = ?", id).First(&channel)
	if result.Error != nil {
		return nil, result.Error
	}
	return &channel, nil
}

func (cp *channelPersistence) GetAllInGuild(guildID *string) ([]*model.Channel, error) {
	var channels []*model.Channel
	result := cp.db.Where("guild_id = ?", *guildID).Find(&channels)
	if result.Error != nil {
		return nil, result.Error
	}
	return channels, nil
}
