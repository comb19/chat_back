package model

import "time"

type UserChannels struct {
	UserID    string
	ChannelID string
	CreatedAt time.Time
	UpdatedAt time.Time
}
