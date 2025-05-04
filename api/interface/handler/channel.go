package handler

import (
	"chat_back/usecase"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type channelURI struct {
	ID string `uri:"channel_id" binding:"required, uuid"`
}

type guildURI struct {
	ID string `uri:"guild_id" binding:"required, uuid"`
}

type Channel struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Private     bool    `json:"private"`
	GuildID     *string `json:"guild_id"`
}

type ChannelHandler interface {
	HandleInsert(ctx *gin.Context)
	HandleGetByID(ctx *gin.Context)
	HandleGetAllInGuild(ctx *gin.Context)
}

type channelHandler struct {
	db             *gorm.DB
	channelUseCase usecase.ChannelUsecase
}

func NewChannelHandler(db *gorm.DB, channelUseCase usecase.ChannelUsecase) ChannelHandler {
	return &channelHandler{
		db:             db,
		channelUseCase: channelUseCase,
	}
}

func (ch *channelHandler) HandleInsert(ctx *gin.Context) {
	var channel Channel

	if err := ctx.BindJSON(&channel); err != nil {
		ctx.String(http.StatusBadRequest, "Invalid request")
		return
	}

	channelId, err := ch.channelUseCase.Insert(ch.db, channel.Name, channel.Description, channel.Private, channel.GuildID)
	if err != nil {
		fmt.Println(err)
		ctx.String(http.StatusInternalServerError, "Failed to insert channel")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"channel_id": *channelId,
	})
}

func (ch *channelHandler) HandleGetByID(ctx *gin.Context) {
	var uri channelURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.String(http.StatusBadRequest, "Invalid request")
		return
	}

	channel, err := ch.channelUseCase.GetByID(ch.db, uri.ID)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to get channel")
		return
	}

	ctx.JSON(http.StatusOK, channel)
}

func (ch *channelHandler) HandleGetAllInGuild(ctx *gin.Context) {
	var uri guildURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.String(http.StatusBadRequest, "Invalid request")
		return
	}

	channels, err := ch.channelUseCase.GetAllInGuild(ch.db, &uri.ID)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to get channels")
		return
	}

	ctx.JSON(http.StatusOK, channels)
}
