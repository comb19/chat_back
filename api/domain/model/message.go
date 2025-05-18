package model

type Message struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Content   string
	UserID    string
	UserName  string
	ChannelID string
	CreatedAt string `gorm:"default"`
	UpdatedAt string
}
