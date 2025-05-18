package types

type GuildURI struct {
	ID string `uri:"guildID" binding:"required,uuid"`
}

type RequestGuild struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ResponseGuild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
