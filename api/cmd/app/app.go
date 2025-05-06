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
	clerkUser "github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	svix "github.com/svix/svix-webhooks/go"
)

type EnvVar struct {
	ClerkSecKey string `env:"CLERK_SECRET_KEY"`
	FrontendUrl string `env:"FRONTEND_URL"`
	SvixSecKey  string `env:"SVIX_SECRET_KEY"`
}

func authenticationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionToken := strings.TrimPrefix(ctx.Request.Header.Get("Authorization"), "Bearer ")
		claims, err := jwt.Verify(ctx.Request.Context(), &jwt.VerifyParams{
			Token: sessionToken,
		})
		if err != nil {
			fmt.Println("unauthorized")
			fmt.Println(err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		user, err := clerkUser.Get(ctx.Request.Context(), claims.Subject)
		if err != nil {
			fmt.Println("user not found")
			fmt.Println(err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		fmt.Println("authorized:", user.Username)
		fmt.Println(user)

		ctx.Set("user", user)
		ctx.Next()
	}
}

func Run() {
	var env_var EnvVar
	if err := env.Parse(&env_var); err != nil {
		fmt.Println(err)
	}

	clerk.SetKey(env_var.ClerkSecKey)

	wh, err := svix.NewWebhook(env_var.SvixSecKey)
	if err != nil {
		fmt.Println("Error creating webhook:", err)
		panic(err)
	}

	db := config.Init()

	userPersistence := persistence.NewUserPersistence()
	userUseCase := usecase.NewUserUsecase(userPersistence)
	userHandler := handler.NewUserHandler(db, wh, userUseCase)

	userChannelsPersistence := persistence.NewUserChannelsPersistence()
	authorizationUseCase := usecase.NewAuthorizationUsecase(userChannelsPersistence)

	messagePersistence := persistence.NewMessagePersistence()
	messageUseCase := usecase.NewMessageUsecase(messagePersistence)
	messageHandler := handler.NewMessageHandler(db, messageUseCase, authorizationUseCase)

	channelPersistence := persistence.NewChannelPersistence()
	channelUseCase := usecase.NewChannelUsecase(userChannelsPersistence, channelPersistence)
	channelHandler := handler.NewChannelHandler(db, channelUseCase)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{env_var.FrontendUrl},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: false,
	}))

	authorized := router.Group("/")
	authorized.Use(authenticationMiddleware())
	{
		authorized.GET("/messages/:channelID", messageHandler.HandleMessageInChannel)

		authorized.GET("/channels", func(ctx *gin.Context) {})
		authorized.GET("/channels/:channelID", func(ctx *gin.Context) {})
		authorized.POST("/channels", channelHandler.HandleInsert)
		authorized.PUT("/channels/:channelID", func(ctx *gin.Context) {})
		authorized.DELETE("/channels/:channelID", func(ctx *gin.Context) {})
	}

	router.POST("/users", userHandler.HandleCreateUserByClerk)

	router.GET("/ws/messages/:channelID", messageHandler.HandleMessageWebSocket)

	router.GET("/ping", func(ctx *gin.Context) {
		fmt.Println("pong")
		ctx.String(http.StatusOK, "pong")
	})

	router.Run("0.0.0.0:8080")
}
