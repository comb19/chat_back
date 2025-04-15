package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type db_env struct {
	user string
	pass string
	port string
}

func Init() *gorm.DB {
	var dbEnv db_env
	if err := env.Parse(&dbEnv); err != nil {
		fmt.Print(err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=disable TimeZone=Asia/Tokyo", "localhost", dbEnv.user, dbEnv.pass, dbEnv.port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	return db
}
