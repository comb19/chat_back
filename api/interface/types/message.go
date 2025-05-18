package types

type Message struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	ChannelID string `json:"channel_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type MessageURI struct {
	ChannelID string `uri:"channelID" binding:"required,uuid"`
}
