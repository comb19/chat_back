package handler

import (
	"chat_back/interface/types"
	"chat_back/usecase"
	"log/slog"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-gonic/gin"
)

type ChannelHandler interface {
	HandleInsert(ctx *gin.Context)
	HandleDelete(ctx *gin.Context)
	HandleGetByID(ctx *gin.Context)
	HandleAddUserToChannel(ctx *gin.Context)
	HandleGetMessagesInChannel(ctx *gin.Context)
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

	var channel types.RequestChannel
	if err := ctx.BindJSON(&channel); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return
	}

	tempUser, ok := ctx.Get("user")
	if !ok {
		ctx.Status(http.StatusUnauthorized)
		return
	}
	user, ok := tempUser.(*clerk.User)
	if !ok {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	newChannel, err := ch.channelUseCase.Insert(channel.Name, channel.Description, channel.Private, user.ID, channel.GuildID)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, types.ResponseChannel{
		ID:          newChannel.ID,
		Name:        newChannel.Name,
		Description: newChannel.Description,
		Private:     newChannel.Private,
		GuildID:     newChannel.GuildID,
	})
}

func (ch *channelHandler) HandleDelete(ctx *gin.Context) {
	slog.DebugContext(ctx, "HandleDelete")

	var uri types.ChannelURI
	if err := ctx.BindUri(&uri); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return
	}
	err := ch.channelUseCase.Delete(uri.ID)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (ch *channelHandler) HandleGetByID(ctx *gin.Context) {
	slog.DebugContext(ctx, "HandleGetByID")

	var uri types.ChannelURI
	if err := ctx.BindUri(&uri); err != nil {
		return
	}

	channel, err := ch.channelUseCase.GetByID(uri.ID)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, types.ResponseChannel{
		ID:          channel.ID,
		Name:        channel.Name,
		Description: channel.Description,
		Private:     channel.Private,
		GuildID:     channel.GuildID,
	})
}

func (ch *channelHandler) HandleAddUserToChannel(ctx *gin.Context) {
	slog.DebugContext(ctx, "HandleAddUserToChannel")
	slog.DebugContext(ctx, ctx.Request.URL.RawPath)

	var uri types.ChannelURI
	if err := ctx.BindUri(&uri); err != nil {
		slog.ErrorContext(ctx, err.Error())
		ctx.Status(http.StatusBadRequest)
		return
	}
	var userRequestBody types.UsersRequestBody
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

func (ch *channelHandler) HandleGetMessagesInChannel(ctx *gin.Context) {
	slog.DebugContext(ctx, "HandleGetMessagesInChannel")

	var uri types.ChannelURI
	if err := ctx.BindUri(&uri); err != nil {
		slog.Error(err.Error())
		return
	}
	slog.DebugContext(ctx, uri.ID)

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

	messages, err := ch.channelUseCase.GetMessagesOfChannel(uri.ID, user.ID)
	if err != nil {
		slog.Error(err.Error())
		ctx.Status(http.StatusInternalServerError)
		return
	}

	responseMessages := make([]types.Message, len(*messages))
	for index, message := range *messages {
		responseMessages[index] = types.Message{
			ID:        message.ID,
			UserID:    message.ID,
			UserName:  message.UserName,
			ChannelID: message.ChannelID,
			Content:   message.Content,
			CreatedAt: message.CreatedAt,
		}
	}
	ctx.JSON(http.StatusOK, responseMessages)
}
