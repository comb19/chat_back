package handler

import (
	"chat_back/interface/types"
	"chat_back/usecase"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-gonic/gin"
)

type GuildInvitationHandler interface {
	CreateGuildInvitation(ctx *gin.Context)
	VerifyGuildInvitation(ctx *gin.Context)
}

type guildInvitationHandler struct {
	guildInvitationUsecase usecase.GuildInvitationUsecase
}

func NewGuildInvitationHandler(guildInvitationUsecase usecase.GuildInvitationUsecase) GuildInvitationHandler {
	return &guildInvitationHandler{
		guildInvitationUsecase: guildInvitationUsecase,
	}
}

func (gih guildInvitationHandler) CreateGuildInvitation(ctx *gin.Context) {
	slog.DebugContext(ctx, "CreateGuildInvitation")

	var requestGuildInvitation types.RequestGuildInvitation
	if err := ctx.BindJSON(&requestGuildInvitation); err != nil {
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

	guildInvitation, err := gih.guildInvitationUsecase.CreateGuildInvitation(user.ID, requestGuildInvitation.GuildID)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, types.ResponseGuildInvitation{
		ID:         guildInvitation.ID,
		OwnerID:    guildInvitation.OwnerID,
		GuildID:    guildInvitation.GuildID,
		Expiration: guildInvitation.Expiration,
		URL:        fmt.Sprintf("/invitations/guilds/%s", guildInvitation.ID),
	})
}

func (gih guildInvitationHandler) VerifyGuildInvitation(ctx *gin.Context) {
	slog.DebugContext(ctx, "VerifyGuildInvitation")

	var guildInvitationUri types.GuildInvitationURI
	if err := ctx.BindUri(&guildInvitationUri); err != nil {
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

	verified, err := gih.guildInvitationUsecase.VerifyGuildInvitation(guildInvitationUri.ID, user.ID)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	if verified {
		ctx.Status(http.StatusAccepted)
		return
	} else {
		ctx.Status(http.StatusBadRequest)
		return
	}
}
