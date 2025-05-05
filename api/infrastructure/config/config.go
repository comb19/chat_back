package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type db_env struct {
	DSN string `env:"DATABASE_URL"`
}

func Connect(dsn string, timeout, interval time.Duration) *gorm.DB {
	deadline := time.Now().Add(timeout)

	for {
		if time.Now().After(deadline) {
			panic("failed to connect to database")
		}
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			fmt.Println("Connected to database")
			return db
		}
		fmt.Println("Retrying to connect to database...")
		time.Sleep(interval)
	}
}

func Init() *gorm.DB {
	var dbEnv db_env
	if err := env.Parse(&dbEnv); err != nil {
		fmt.Print(err)
	}
	return Connect(dbEnv.DSN, 3*time.Minute, 15*time.Second)
}
