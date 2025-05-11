package persistence

import (
	"chat_back/domain/model"
	"chat_back/domain/repository"
	"errors"

	"gorm.io/gorm"
)

type guildPersistence struct {
	db *gorm.DB
}

func NewGuildPersistence(db *gorm.DB) repository.GuildRepository {
	return &guildPersistence{
		db: db,
	}
}

func (gp guildPersistence) Insert(name, description string) (*model.Guild, error) {
	guild := &model.Guild{
		Name:        name,
		Description: description,
	}
	result := gp.db.Select("name", "description").Create(guild)
	if result.Error != nil {
		return nil, result.Error
	}
	return guild, nil
}

func (gp guildPersistence) Find(id string) (*model.Guild, error) {
	var guild *model.Guild
	result := gp.db.Where("id = ?", id).First(guild)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return guild, nil
}

func (gp guildPersistence) FindOfUser(userID string) ([]*model.Guild, error) {
	var guilds []*model.Guild
	result := gp.db.Model(&model.Guild{}).Where("user_guilds.user_id = ?", userID).Joins("inner join user_guilds on guilds.id = user_guilds.guild_id").Scan(guilds)
	if result.Error != nil {
		return []*model.Guild{}, result.Error
	}
	return guilds, nil
}
