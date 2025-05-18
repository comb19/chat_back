package types

type ChannelURI struct {
	ID string `uri:"channelID" binding:"required,uuid"`
}

type UsersRequestBody struct {
	UserIDs []string `json:"user_ids" binding:"required"`
}

type RequestChannel struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Private     bool    `json:"private"`
	GuildID     *string `json:"guild_id"`
}

type ResponseChannel struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Private     bool    `json:"private"`
	GuildID     *string `json:"guild_id"`
}
