package handler

import (
	"fmt"
	"net/http"
	"todo_back/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type todo struct {
	Title       string
	Description string
}

type TodoHandler interface {
	HandleTodoInsert(ctx *gin.Context)
	HandleTodoGetAll(ctx *gin.Context)
	// HandleTodoGetByID(ctx *gin.Context)
	// HandleTodoUpdateByID(ctx *gin.Context)
	// HandleTodoDeleteByID(ctx *gin.Context)
}

type todoHandler struct {
	db          *gorm.DB
	todoUseCase usecase.TodoUsecase
}

func (th todoHandler) HandleTodoInsert(ctx *gin.Context) {
	var todo todo

	if err := ctx.BindJSON(&todo); err != nil {
		fmt.Println(err)
		return
	}
	if err := th.todoUseCase.Insert(th.db, todo.Title, todo.Description); err != nil {
		fmt.Println(err)
		ctx.String(http.StatusInternalServerError, "Failed to insert todo")
		return
	}
	ctx.String(http.StatusOK, "Inserted")
}

func (th todoHandler) HandleTodoGetAll(ctx *gin.Context) {
	todos, err := th.todoUseCase.GetAll(th.db)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to get todos")
		return
	}
	ctx.JSON(http.StatusOK, todos)
}

func NewTodoHandler(db *gorm.DB, todoUseCase usecase.TodoUsecase) TodoHandler {
	return &todoHandler{
		db:          db,
		todoUseCase: todoUseCase,
	}
}
