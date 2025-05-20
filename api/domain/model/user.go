package model

import "time"

type User struct {
	ID        string
	UserName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
