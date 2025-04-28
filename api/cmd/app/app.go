package app

import (
	"fmt"
	"net/http"
	"strings"
	"todo_back/infrastructure/config"
	"todo_back/infrastructure/persistence"
	"todo_back/interface/handler"
	"todo_back/usecase"

	"github.com/caarlos0/env/v11"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type EnvVar struct {
	Clerk_sec_key string `env:"CLERK_SECRET_KEY"`
}

func authenticationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionToken := strings.TrimPrefix(ctx.Request.Header.Get("Authorization"), "Bearer ")
		fmt.Println(sessionToken)
		claims, err := jwt.Verify(ctx.Request.Context(), &jwt.VerifyParams{
			Token: sessionToken,
		})
		if err != nil {
			fmt.Println("unauthorized")
			fmt.Println(err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		_, err = user.Get(ctx.Request.Context(), claims.Subject)
		if err != nil {
			fmt.Println("user not found")
			fmt.Println(err)
			ctx.JSON((http.StatusNotFound), gin.H{"error": "user not found"})
			return
		}
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

	clerk.SetKey(env_var.Clerk_sec_key)
	// clerk_config := &clerk.ClientConfig{}
	// clerk_config.Key = &env_var.clerk_sec_key
	// clerk_client := user.NewClient(clerk_config)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: false,
	}))

	// router.Use(authenticationMiddleware(clerk_client))
	router.Use(authenticationMiddleware())

	router.POST("/todos", todoHandler.HandleTodoInsert)
	router.GET("/todos", todoHandler.HandleTodoGetAll)

	router.GET("/ping", func(ctx *gin.Context) {
		fmt.Println("pong")
		ctx.String(http.StatusOK, "pong")
	})

	router.Run("0.0.0.0:8080")
}
