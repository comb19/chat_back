package repository

type GuildRepository interface {
	Insert(guildID string, name string) error
	GetByID(guildID string) (string, error)
}
