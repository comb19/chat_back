package repository

import (
	"chat_back/domain/model"
	"time"
)

type GuildInvitationRepository interface {
	Insert(ownerID, guildID string, expiration time.Time) (*model.GuildInvitation, error)
	Find(id string) (*model.GuildInvitation, error)
}
