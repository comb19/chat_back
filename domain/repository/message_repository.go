package repository

type MessageRepository interface {
	Insert(guildID string, channelID string, userID string, content string) error
	GetByID(guildID string, channelID string, userID string) (string, error)
	GetAllInChannel(guildID string, channelID string) ([]string, error)
}
