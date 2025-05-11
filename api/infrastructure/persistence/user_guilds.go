package persistence

import (
	"chat_back/domain/model"
	"chat_back/domain/repository"
	"errors"

	"gorm.io/gorm"
)

type userGuildsPersistence struct {
	db *gorm.DB
}

func NewUserGuildsPersistence(db *gorm.DB) repository.UserGuildsRepository {
	return &userGuildsPersistence{
		db: db,
	}
}

func (ug userGuildsPersistence) Insert(userID string, guildID string) (*model.UserGuilds, error) {
	userGuild := model.UserGuilds{
		UserID:  userID,
		GuildID: guildID,
	}
	result := ug.db.Select("user_id", "guild_id").Create(&userGuild)
	if result.Error != nil {
		return nil, result.Error
	}

	return &userGuild, nil
}

func (ug userGuildsPersistence) Find(userID string, guildID string) (*model.UserGuilds, error) {
	var userGuilds *model.UserGuilds
	result := ug.db.Where("user_id = ? AND guild_id = ?", userID, guildID).First(userGuilds)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return userGuilds, nil
}
