package repository

import "gorm.io/gorm"

type UserRepository interface {
	Insert(DB *gorm.DB, name string) error
	GetByID(DB *gorm.DB, id string) (string, error)
}
