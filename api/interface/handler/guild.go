package handler

import (
	"chat_back/interface/types"
	"chat_back/usecase"
	"log/slog"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-gonic/gin"
)

type GuildHandler interface {
	HandlePostGuilds(ctx *gin.Context)
	HandleGetGuilds(ctx *gin.Context)
	HandleGetChannelsOfGuild(ctx *gin.Context)
	HandleCreateChannelInGuild(ctx *gin.Context)
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

	var requestGuild types.RequestGuild
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

	ctx.JSON(http.StatusCreated, types.ResponseGuild{
		ID:          guild.ID,
		Name:        guild.Name,
		Description: guild.Description,
	})
}

func (gh guildHandler) HandleGetGuilds(ctx *gin.Context) {
	slog.DebugContext(ctx, "HandleGetGuilds")

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

	guilds, err := gh.guildUseCase.GetGuildsOfUser(user.ID)
	if err != nil {
		slog.Error(err.Error())
		ctx.Status(http.StatusInternalServerError)
		return
	}

	responseGuilds := make([]types.ResponseGuild, len(guilds))
	for index, guild := range guilds {
		responseGuilds[index] = types.ResponseGuild{
			ID:          guild.ID,
			Name:        guild.Name,
			Description: guild.Description,
		}
	}

	ctx.JSON(http.StatusOK, responseGuilds)
}

func (gh guildHandler) HandleGetChannelsOfGuild(ctx *gin.Context) {
	slog.DebugContext(ctx, "HandleGetChannelsOfGuild")

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

	var uri types.GuildURI
	if err := ctx.BindUri(&uri); err != nil {
		return
	}

	channels, err := gh.guildUseCase.GetChannelsOfGuild(uri.ID, user.ID)
	if err != nil {
		slog.Error(err.Error())
		ctx.Status(http.StatusInternalServerError)
		return
	}

	responseChannels := make([]*types.ResponseChannel, len(channels))
	for index, channel := range channels {
		responseChannels[index] = &types.ResponseChannel{
			ID:          channel.ID,
			Name:        channel.Name,
			Description: channel.Description,
			Private:     channel.Private,
		}
	}

	ctx.JSON(http.StatusOK, responseChannels)
}

func (gh guildHandler) HandleCreateChannelInGuild(ctx *gin.Context) {
	slog.DebugContext(ctx, "HandleCreateChannelInGuild")

	var requestChannel types.RequestChannel
	if err := ctx.BindJSON(&requestChannel); err != nil {
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

	var uri types.GuildURI
	if err := ctx.BindUri(&uri); err != nil {
		return
	}

	channel, err := gh.guildUseCase.CreateNewChannelInGuild(uri.ID, user.ID, requestChannel.Name, requestChannel.Description)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, channel)
}
