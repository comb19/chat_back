package handler

import (
	"chat_back/usecase"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-gonic/gin"
)

type usersRequestBody struct {
	UserIDs []string `json:"user_ids" binding:"required"`
}

type channelURI struct {
	ID string `uri:"channelID" binding:"required,uuid"`
}

type guildURI struct {
	ID string `uri:"guildID" binding:"required,uuid"`
}

type RequestChannel struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Private     bool    `json:"private"`
	GuildID     *string `json:"guild_id"`
}

type ResponseChannel struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Private     bool    `json:"private"`
	GuildID     *string `json:"guild_id"`
}

type ChannelHandler interface {
	HandleInsert(ctx *gin.Context)
	HandleGetByID(ctx *gin.Context)
	HandleGetAllInGuild(ctx *gin.Context)
	HandleAddUserToChannel(ctx *gin.Context)
}

type channelHandler struct {
	channelUseCase usecase.ChannelUsecase
}

func NewChannelHandler(channelUseCase usecase.ChannelUsecase) ChannelHandler {
	return &channelHandler{
		channelUseCase: channelUseCase,
	}
}

func (ch *channelHandler) HandleInsert(ctx *gin.Context) {
	slog.DebugContext(ctx, "HandleInsert")

	var channel RequestChannel
	if err := ctx.BindJSON(&channel); err != nil {
		slog.ErrorContext(ctx, err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	tempUser, ok := ctx.Get("user")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, ok := tempUser.(*clerk.User)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	newChannel, err := ch.channelUseCase.Insert(channel.Name, channel.Description, channel.Private, user.ID, channel.GuildID)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert channel"})
		return
	}

	ctx.JSON(http.StatusCreated, ResponseChannel{
		ID:          newChannel.ID,
		Name:        newChannel.Name,
		Description: newChannel.Description,
		Private:     newChannel.Private,
		GuildID:     newChannel.GuildID,
	})
}

func (ch *channelHandler) HandleGetByID(ctx *gin.Context) {
	slog.DebugContext(ctx, "HandleGetByID")

	var uri channelURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	channel, err := ch.channelUseCase.GetByID(uri.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get channel"})
		return
	}

	ctx.JSON(http.StatusOK, ResponseChannel{
		ID:          channel.ID,
		Name:        channel.Name,
		Description: channel.Description,
		Private:     channel.Private,
		GuildID:     channel.GuildID,
	})
}

func (ch *channelHandler) HandleGetAllInGuild(ctx *gin.Context) {
	slog.DebugContext(ctx, "HandleGetAllInGuild")

	var uri guildURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	channels, err := ch.channelUseCase.GetAllInGuild(&uri.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get channels"})
		return
	}

	responseChannels := make([]ResponseChannel, 0)
	for _, channel := range channels {
		responseChannels = append(responseChannels, ResponseChannel{
			ID:          channel.ID,
			Name:        channel.Name,
			Description: channel.Description,
			Private:     channel.Private,
			GuildID:     channel.GuildID,
		})
	}

	ctx.JSON(http.StatusOK, responseChannels)
}

func (ch *channelHandler) HandleAddUserToChannel(ctx *gin.Context) {
	slog.DebugContext(ctx, "HandleAddUserToChannel")
	slog.DebugContext(ctx, ctx.Request.URL.RawPath)

	var uri channelURI
	if err := ctx.BindUri(&uri); err != nil {
		slog.ErrorContext(ctx, err.Error())
		ctx.Status(http.StatusBadRequest)
		return
	}
	var userRequestBody usersRequestBody
	if err := ctx.BindJSON(&userRequestBody); err != nil {
		return
	}

	_, err := ch.channelUseCase.AddUserToChannel(uri.ID, userRequestBody.UserIDs)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.Status(http.StatusCreated)
}
