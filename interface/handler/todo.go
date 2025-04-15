package handler

import (
	"net/http"
	"todo_back/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type todo struct {
	title       string
	description string
}

type TodoHandler interface {
	HandleTodoInsert(ctx *gin.Context)
}

type todoHandler struct {
	db          *gorm.DB
	todoUseCase usecase.TodoUsecase
}

func (th todoHandler) HandleTodoInsert(ctx *gin.Context) {
	var todo todo

	if err := ctx.ShouldBind(&todo); err != nil {
		th.todoUseCase.Insert(th.db, todo.title, todo.description)
		ctx.String(http.StatusOK, "Inserted")
	}
}

func NewTodoHandler(db *gorm.DB, todoUseCase usecase.TodoUsecase) TodoHandler {
	return &todoHandler{
		db:          db,
		todoUseCase: todoUseCase,
	}
}
