package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type db_env struct {
	Host string `env:"POSTGRES_HOSTNAME"`
	User string `env:"POSTGRES_USER"`
	Pass string `env:"POSTGRES_PASSWORD"`
	Port string `env:"DB_PORT"`
	DB   string `env:"POSTGRES_DB"`
}

func Init() *gorm.DB {
	var dbEnv db_env
	if err := env.Parse(&dbEnv); err != nil {
		fmt.Print(err)
	}
	fmt.Println(dbEnv)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo", dbEnv.Host, dbEnv.User, dbEnv.Pass, dbEnv.DB, dbEnv.Port)
	fmt.Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	return db
}
