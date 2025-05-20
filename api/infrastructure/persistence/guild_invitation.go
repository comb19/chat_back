package persistence

import (
	"chat_back/domain/model"
	"chat_back/domain/repository"
	"errors"
	"time"

	"gorm.io/gorm"
)

type guildInvitationPersistence struct {
	db *gorm.DB
}

func NewGuildInvitationPersistence(db *gorm.DB) repository.GuildInviationRepository {
	return &guildInvitationPersistence{
		db: db,
	}
}

func (gip guildInvitationPersistence) Insert(ownerID, guildID string, expiration time.Time) (*model.GuildInvitation, error) {
	guildInvitation := model.GuildInvitation{
		OwnerID:    ownerID,
		GuildID:    guildID,
		Expiration: expiration,
	}
	result := gip.db.Create(&guildInvitation)
	if err := result.Error; err != nil {
		return nil, err
	}

	return &guildInvitation, nil
}

func (gip guildInvitationPersistence) Find(id string) (*model.GuildInvitation, error) {
	var guildInvitation model.GuildInvitation
	result := gip.db.Where("id = ?", id).First(&guildInvitation)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &guildInvitation, nil
}
