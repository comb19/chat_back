package app

import (
	"fmt"
	"todo_back/infrastructure/config"
	"todo_back/infrastructure/persistence"
	"todo_back/interface/handler"
	"todo_back/usecase"

	"github.com/caarlos0/env"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type EnvVar struct {
	clerk_sec_key string `env:"CLERK_SECRET_KEY"`
}

func Run() {
	var env_var EnvVar
	if err := env.Parse(&env_var); err != nil {
		fmt.Println(err)
	}

	db := config.Init()

	todoPersistence := persistence.NewTodoPersistence()
	todoUseCase := usecase.NewTodoUsecase(todoPersistence)
	todoHandler := handler.NewTodoHandler(db, todoUseCase)

	clerk.SetKey(env_var.clerk_sec_key)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: false,
	}))

	router.POST("/todos", todoHandler.HandleTodoInsert)
	router.GET("/todos", todoHandler.HandleTodoGetAll)

	router.Run("0.0.0.0:8080")
}
