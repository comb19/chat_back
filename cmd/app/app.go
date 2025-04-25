package app

import (
	"fmt"
	"net/http"
	"todo_back/infrastructure/config"
	"todo_back/infrastructure/persistence"
	"todo_back/interface/handler"
	"todo_back/usecase"

	"github.com/caarlos0/env/v11"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type EnvVar struct {
	clerk_sec_key string `env:"CLERK_SECRET_KEY"`
}

func authenticationMiddleware(client *user.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, ok := clerk.SessionClaimsFromContext(ctx.Request.Context())
		if !ok {
			fmt.Print("unauthorized")
			return
		}
		usr, err := user.Get(ctx.Request.Context(), claims.Subject)
		if err != nil {
			fmt.Println("error")
		}
		fmt.Printf("%s %s\n", usr.ID, *usr.FirstName)
	}
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

	clerk_config := &clerk.ClientConfig{}
	clerk_config.Key = &env_var.clerk_sec_key
	clerk_client := user.NewClient(clerk_config)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: false,
	}))

	router.Use(authenticationMiddleware(clerk_client))

	router.POST("/todos", todoHandler.HandleTodoInsert)
	router.GET("/todos", todoHandler.HandleTodoGetAll)

	router.GET("ping", func(ctx *gin.Context) {
		fmt.Println("pong")
		ctx.String(http.StatusOK, "pong")
	})

	router.Run("0.0.0.0:8080")
}
