package model

import "time"

type GuildInvitation struct {
	ID         string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OwnerID    string
	GuildID    string
	Expiration time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
