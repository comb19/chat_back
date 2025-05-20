package model

import "time"

type UserGuilds struct {
	UserID    string
	GuildID   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
