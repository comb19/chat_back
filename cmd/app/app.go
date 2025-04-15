package app

import (
	"fmt"
	"os"
	"todo_back/infrastructure/config"
	"todo_back/infrastructure/persistence"
	"todo_back/interface/handler"
	"todo_back/usecase"

	"github.com/caarlos0/env/v11"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type test struct {
	Port string `env:"DB_PORT"`
}

func Run() {
	var t test
	if err := env.Parse(&t); err != nil {
		fmt.Println(err)
	}
	fmt.Println(t.Port)
	fmt.Println(os.Getenv("DB_PORT"))
	fmt.Println("hello")
	db := config.Init()

	todoPersistence := persistence.NewTodoPersistence()
	todoUseCase := usecase.NewTodoUsecase(todoPersistence)
	todoHandler := handler.NewTodoHandler(db, todoUseCase)

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
