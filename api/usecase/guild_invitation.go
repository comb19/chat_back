package usecase

import (
	"chat_back/domain/model"
	"chat_back/domain/repository"
	"fmt"
	"log/slog"
	"time"
)

const expireDays = 1

type GuildInvitationUsecase interface {
	CreateGuildInvitation(ownerID, guildID string) (*model.GuildInvitation, error)
	VerifyGuildInvitation(id, userID string) (bool, *model.GuildInvitation, error)
}

type guildInvitationUsecase struct {
	guildInvitationRepository repository.GuildInvitationRepository
	userGuildRepository       repository.UserGuildsRepository
}

func NewGuildInvitationUsecase(guildInvitationRepository repository.GuildInvitationRepository, userGuildRepository repository.UserGuildsRepository) GuildInvitationUsecase {
	return guildInvitationUsecase{
		guildInvitationRepository: guildInvitationRepository,
		userGuildRepository:       userGuildRepository,
	}
}

func (giu guildInvitationUsecase) CreateGuildInvitation(ownerID, guildID string) (*model.GuildInvitation, error) {
	return giu.guildInvitationRepository.Insert(ownerID, guildID, time.Now())
}

func (giu guildInvitationUsecase) VerifyGuildInvitation(id, userID string) (bool, *model.GuildInvitation, error) {
	guildInvitation, err := giu.guildInvitationRepository.Find(id)
	if err != nil {
		slog.Error(err.Error())
		return false, nil, err
	}
	if guildInvitation == nil {
		slog.Debug(fmt.Sprintf("the guild invitation %s is not found.", id))
		return false, nil, nil
	}
	if !guildInvitation.Expiration.Before(time.Now().AddDate(0, 0, expireDays)) {
		slog.Debug(fmt.Sprintf("the guild invitation %s is expired.\n", id))
		slog.Debug(fmt.Sprintf("Expire: %s, Now: %s.\n", guildInvitation.Expiration, time.Now()))
		return false, nil, nil
	}

	userGuild, err := giu.userGuildRepository.Find(userID, guildInvitation.GuildID)
	if err != nil {
		slog.Error(err.Error())
		return false, nil, err
	}
	if userGuild != nil {
		slog.Debug(fmt.Sprintf("The user %s already belongs to the guild %s.", userID, guildInvitation.GuildID))
		return true, guildInvitation, err
	}

	_, err = giu.userGuildRepository.Insert(userID, guildInvitation.GuildID)
	if err != nil {
		slog.Error(err.Error())
		return false, nil, err
	}
	slog.Debug(fmt.Sprintf("The user %s joined the guild %s.", userID, guildInvitation.GuildID))
	return true, guildInvitation, nil
}
