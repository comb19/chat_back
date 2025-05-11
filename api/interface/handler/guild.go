package handler

import (
	"chat_back/usecase"
	"log/slog"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-gonic/gin"
)

type RequestGuild struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ResponseGuild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GuildHandler interface {
	HandlePostGuilds(ctx *gin.Context)
}

type guildHandler struct {
	guildUseCase usecase.GuildUseCase
}

func NewGuildHandler(guildUseCase usecase.GuildUseCase) GuildHandler {
	return &guildHandler{
		guildUseCase: guildUseCase,
	}
}

func (gh guildHandler) HandlePostGuilds(ctx *gin.Context) {
	slog.DebugContext(ctx, "HandlePostGuilds")

	var requestGuild RequestGuild
	if err := ctx.BindJSON(&requestGuild); err != nil {
		return
	}

	tempUser, ok := ctx.Get("user")
	if !ok {
		ctx.Status(http.StatusUnauthorized)
		return
	}
	user, ok := tempUser.(*clerk.User)
	if !ok {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	guild, err := gh.guildUseCase.CreateNewGuild(requestGuild.Name, requestGuild.Description, user.ID)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, ResponseGuild{
		ID:          guild.ID,
		Name:        guild.Name,
		Description: guild.Description,
	})
}
