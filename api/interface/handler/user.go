package handler

import (
	"chat_back/usecase"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	svix "github.com/svix/svix-webhooks/go"
)

type UserHandler interface {
	HandleCreateUserByClerk(ctx *gin.Context)
	HandleGetUserByID(ctx *gin.Context)
}

type userHandler struct {
	wh          *svix.Webhook
	userUseCase usecase.UserUsecase
}

func NewUserHandler(wh *svix.Webhook, userUseCase usecase.UserUsecase) UserHandler {
	return &userHandler{
		wh:          wh,
		userUseCase: userUseCase,
	}
}

type user struct {
	Data struct {
		ID       string `json:"id"`
		UserName string `json:"username"`
	} `json:"data"`
}

func (uh *userHandler) HandleCreateUserByClerk(ctx *gin.Context) {
	headers := ctx.Request.Header
	payload, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		ctx.String(http.StatusBadRequest, "Bad Request")
		return
	}

	if err := uh.wh.Verify(payload, headers); err != nil {
		fmt.Println("Error verifying webhook:", err)
		ctx.String(http.StatusBadRequest, "Bad Request")
		return
	}

	fmt.Println("Webhook verified successfully")
	fmt.Println(headers)
	fmt.Println(string(payload))

	var user user
	if err := json.Unmarshal(payload, &user); err != nil {
		fmt.Println("Error binding JSON:", err)
		ctx.String(http.StatusBadRequest, "Bad Request")
		return
	}
	fmt.Println(user)

	if _, err := uh.userUseCase.CreateUserByClerk(user.Data.ID, user.Data.UserName); err != nil {
		fmt.Println("Error creating user:", err)
		ctx.String(http.StatusInternalServerError, "Failed to create user")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (uh *userHandler) HandleGetUserByID(ctx *gin.Context) {
}
