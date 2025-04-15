package app

import (
	"todo_back/infrastructure/config"
	"todo_back/infrastructure/persistence"
	"todo_back/interface/handler"
	"todo_back/usecase"

	"github.com/gin-gonic/gin"
)

func Run() {
	db := config.Init()

	todoPersistence := persistence.NewTodoPersistence()
	todoUseCase := usecase.NewTodoUsecase(todoPersistence)
	todoHandler := handler.NewTodoHandler(db, todoUseCase)

	router := gin.Default()

	router.GET("/todos", todoHandler.HandleTodoInsert)

	router.Run("localhost:8080")
}
