package model

type Channel struct {
	ID          string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string
	Description string
	GuildID     *string
	Private     bool
	CreatedAt   string
	UpdatedAt   string
}
