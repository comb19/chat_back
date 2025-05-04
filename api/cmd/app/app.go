package app

import (
	"chat_back/infrastructure/config"
	"chat_back/infrastructure/persistence"
	"chat_back/interface/handler"
	"chat_back/usecase"
	"fmt"
	"net/http"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type EnvVar struct {
	Clerk_sec_key string `env:"CLERK_SECRET_KEY"`
	Frontend_url  string `env:"FRONTEND_URL"`
}

func authenticationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionToken := strings.TrimPrefix(ctx.Request.Header.Get("Authorization"), "Bearer ")
		fmt.Println("session token", sessionToken)
		claims, err := jwt.Verify(ctx.Request.Context(), &jwt.VerifyParams{
			Token: sessionToken,
		})
		if err != nil {
			fmt.Println("unauthorized")
			fmt.Println(err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		clerk_user, err := user.Get(ctx.Request.Context(), claims.Subject)
		if err != nil {
			fmt.Println("user not found")
			fmt.Println(err)
			ctx.JSON((http.StatusNotFound), gin.H{"error": "user not found"})
			return
		}
		fmt.Println(clerk_user)
		ctx.Set("user", clerk_user)
		ctx.Next()
	}
}

func Run() {
	var env_var EnvVar
	if err := env.Parse(&env_var); err != nil {
		fmt.Println(err)
	}
	fmt.Println("clerk sec key")
	fmt.Println(env_var.Clerk_sec_key)

	db := config.Init()

	todoPersistence := persistence.NewTodoPersistence()
	todoUseCase := usecase.NewTodoUsecase(todoPersistence)
	todoHandler := handler.NewTodoHandler(db, todoUseCase)

	userChannelsPersistence := persistence.NewUserChannelsPersistence()
	authorizationUseCase := usecase.NewAuthorizationUsecase(userChannelsPersistence)

	messagePersistence := persistence.NewMessagePersistence()
	messageUseCase := usecase.NewMessageUsecase(messagePersistence)
	messageHandler := handler.NewMessageHandler(db, messageUseCase, authorizationUseCase)

	channelPersistence := persistence.NewChannelPersistence()
	channelUseCase := usecase.NewChannelUsecase(channelPersistence)
	channelHandler := handler.NewChannelHandler(db, channelUseCase)

	clerk.SetKey(env_var.Clerk_sec_key)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{env_var.Frontend_url},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: false,
	}))

	authorized := router.Group("/")
	authorized.Use(authenticationMiddleware())
	{
		authorized.POST("/todos", todoHandler.HandleTodoInsert)
		authorized.GET("/todos", todoHandler.HandleTodoGetAll)

		authorized.GET("/messages/:channelID", messageHandler.HandleMessageInChannel)

		authorized.GET("/channels", func(ctx *gin.Context) {})
		authorized.GET("/channels/:channelID", func(ctx *gin.Context) {})
		authorized.POST("/channels", channelHandler.HandleInsert)
		authorized.PUT("/channels/:channelID", func(ctx *gin.Context) {})
		authorized.DELETE("/channels/:channelID", func(ctx *gin.Context) {})
	}

	router.GET("/ws/messages/:channelID", messageHandler.HandleMessageWebSocket)

	router.GET("/ping", func(ctx *gin.Context) {
		fmt.Println("pong")
		ctx.String(http.StatusOK, "pong")
	})

	router.Run("0.0.0.0:8080")
}
