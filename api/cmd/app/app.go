package app

import (
	"chat_back/infrastructure/config"
	"chat_back/infrastructure/persistence"
	"chat_back/interface/handler"
	"chat_back/usecase"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	clerkUser "github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/clog"
	svix "github.com/svix/svix-webhooks/go"
)

type EnvVar struct {
	Development bool   `env:"DEVELOPMENT" envDefault:"false"`
	ClerkSecKey string `env:"CLERK_SECRET_KEY,notEmpty"`
	FrontendUrl string `env:"FRONTEND_URL,notEmpty"`
	SvixSecKey  string `env:"SVIX_SECRET_KEY,notEmpty"`
}

var contextKeys = []string{"user"}

type LogHandler struct {
	slog.Handler
}

func (lh *LogHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, key := range contextKeys {
		if val := ctx.Value(key); val != nil {
			switch key {
			case "user":
				if user, ok := val.(*clerk.User); ok {
					r.AddAttrs(slog.Attr{Key: string(key), Value: slog.StringValue(fmt.Sprintf("id:%s name:%s", user.ID, *user.Username))})
				} else {
					r.AddAttrs(slog.Attr{Key: string(key), Value: slog.AnyValue(val)})
				}
			default:
				r.AddAttrs(slog.Attr{Key: string(key), Value: slog.AnyValue(val)})
			}
		}
	}
	return lh.Handler.Handle(ctx, r)
}

func InitLog(development bool) {
	var slogLevel slog.Level
	if development {
		slogLevel = slog.LevelDebug
	} else {
		slogLevel = slog.LevelInfo
	}
	logger := slog.New(&LogHandler{
		clog.New(
			clog.WithColor(true),
			clog.WithSource(true),
			clog.WithLevel(slogLevel),
		),
	})
	slog.SetDefault(logger)
}

func authenticationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		slog.DebugContext(ctx, "authenticationMiddleware")

		sessionToken := strings.TrimPrefix(ctx.Request.Header.Get("Authorization"), "Bearer ")
		claims, err := jwt.Verify(ctx.Request.Context(), &jwt.VerifyParams{
			Token: sessionToken,
		})
		if err != nil {
			slog.Error("unauthorized")
			slog.Error(err.Error())
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		user, err := clerkUser.Get(ctx.Request.Context(), claims.Subject)
		if err != nil {
			slog.Error("user not found")
			slog.Error(err.Error())
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		ctx.Set("user", user)
		ctx.Next()

		slog.InfoContext(ctx, "authorized")
	}
}

func SetupRouter() *gin.Engine {
	var envVar EnvVar
	if err := env.Parse(&envVar); err != nil {
		panic(err)
	}
	fmt.Println(envVar.Development)

	InitLog(envVar.Development)

	slog.Info("Run")

	clerk.SetKey(envVar.ClerkSecKey)

	wh, err := svix.NewWebhook(envVar.SvixSecKey)
	if err != nil {
		panic(err)
	}

	db := config.Init()

	userPersistence := persistence.NewUserPersistence(db)
	userUseCase := usecase.NewUserUsecase(userPersistence)
	userHandler := handler.NewUserHandler(wh, userUseCase)

	userChannelsPersistence := persistence.NewUserChannelsPersistence(db)
	userGuildsPersistence := persistence.NewUserGuildsPersistence(db)

	authorizationUseCase := usecase.NewAuthorizationUsecase(userChannelsPersistence)

	messagePersistence := persistence.NewMessagePersistence(db)
	messageUseCase := usecase.NewMessageUsecase(messagePersistence)
	messageHandler := handler.NewMessageHandler(messageUseCase, authorizationUseCase)

	channelPersistence := persistence.NewChannelPersistence(db)
	channelUseCase := usecase.NewChannelUsecase(userChannelsPersistence, channelPersistence, messagePersistence)
	channelHandler := handler.NewChannelHandler(channelUseCase)

	guildPersistence := persistence.NewGuildPersistence(db)
	guildUseCase := usecase.NewGuildUseCase(guildPersistence, userGuildsPersistence, channelPersistence, userChannelsPersistence)
	guildHandler := handler.NewGuildHandler(guildUseCase)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{envVar.FrontendUrl},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: false,
	}))

	authorized := router.Group("/")
	authorized.Use(authenticationMiddleware())
	{
		authorized.POST("/channels", channelHandler.HandleInsert)
		authorized.GET("/channels/:channelID", channelHandler.HandleGetByID)
		authorized.PUT("/channels/:channelID", func(ctx *gin.Context) {})
		authorized.DELETE("/channels/:channelID", func(ctx *gin.Context) {})
		authorized.GET("/channels/:channelID/users", func(ctx *gin.Context) {})
		authorized.POST("/channels/:channelID/users", channelHandler.HandleAddUserToChannel)
		authorized.GET("/channels/:channelID/messages", channelHandler.HandleGetMessagesInChannel)

		authorized.GET("/guilds", guildHandler.HandleGetGuilds)
		authorized.POST("/guilds", guildHandler.HandlePostGuilds)
		authorized.GET("/guilds/:guildID", func(ctx *gin.Context) {})
		authorized.PUT("/guilds/:guildID", func(ctx *gin.Context) {})
		authorized.DELETE("/guilds/:guildID", func(ctx *gin.Context) {})
		authorized.GET("/guilds/:guildID/channels", guildHandler.HandleGetChannelsOfGuild)
		authorized.POST("/guilds/:guildID/channels", guildHandler.HandleCreateChannelInGuild)
		authorized.GET("/guilds/:guildID/users", func(ctx *gin.Context) {})
	}

	router.POST("/users", userHandler.HandleCreateUserByClerk)

	router.GET("/ws/messages/:channelID", messageHandler.HandleMessageWebSocket)

	router.GET("/ping", func(ctx *gin.Context) {
		slog.Info("pong")
		ctx.String(http.StatusOK, "pong")
	})

	return router
}

func Run() {
	router := SetupRouter()

	router.Run("0.0.0.0:8080")
}
