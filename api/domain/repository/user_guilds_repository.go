package repository

import (
	"chat_back/domain/model"
)

type UserGuildsRepository interface {
	Insert(userID string, guildID string) (*model.UserGuilds, error)
	Find(userID string, guildID string) (*model.UserGuilds, error)
}
