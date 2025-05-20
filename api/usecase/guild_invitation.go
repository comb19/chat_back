package usecase

import (
	"chat_back/domain/model"
	"chat_back/domain/repository"
	"time"
)

const expireDays = 1

type GuildInvitationUsecase interface {
	CreateGuildInvitation(ownerID, guildID string) (*model.GuildInvitation, error)
	VerifyGuildInvitation(id, userID string) (bool, error)
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

func (giu guildInvitationUsecase) VerifyGuildInvitation(id, userID string) (bool, error) {
	guildInvitation, err := giu.guildInvitationRepository.Find(id)
	if err != nil {
		return false, err
	}
	if guildInvitation == nil {
		return false, nil
	}
	if !guildInvitation.Expiration.Before(time.Now().AddDate(0, 0, expireDays)) {
		return false, nil
	}

	userGuild, err := giu.userGuildRepository.Find(userID, guildInvitation.GuildID)
	if err != nil {
		return false, err
	}
	if userGuild != nil {
		return true, err
	}

	_, err = giu.userGuildRepository.Insert(userID, guildInvitation.GuildID)
	if err != nil {
		return false, err
	}
	return true, nil
}
