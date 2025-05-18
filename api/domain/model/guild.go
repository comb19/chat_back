package model

type Guild struct {
	ID          string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string
	Description string
	CreatedAt   string
	UpdatedAt   string
}
