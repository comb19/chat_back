package types

import "time"

type GuildInvitationURI struct {
	ID string `uri:"invitationID" binding:"required,uuid"`
}

type RequestGuildInvitation struct {
	GuildID string `json:"guild_id"`
}

type ResponseGuildInvitation struct {
	ID         string    `json:"id"`
	OwnerID    string    `json:"owner_id"`
	GuildID    string    `json:"guild_id"`
	Expiration time.Time `json:"expiration"`
	URL        string    `json:"url"`
}
