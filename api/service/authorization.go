package service

import "chat_back/domain/repository"

type AuthorizationService interface {
	CheckAuthorizationAccessToChannel(userID, channelID string) (bool, error)
}

type authorizationService struct {
	userGuildsRepository   repository.UserGuildsRepository
	userChannelsRepository repository.UserChannelsRepository
	channelRepository      repository.ChannelRepository
}

func NewAuthorizationService(userGuildsRepository repository.UserGuildsRepository, userChannelsRepository repository.UserChannelsRepository, channelRepository repository.ChannelRepository) AuthorizationService {
	return &authorizationService{
		userGuildsRepository:   userGuildsRepository,
		userChannelsRepository: userChannelsRepository,
		channelRepository:      channelRepository,
	}
}

func (as authorizationService) CheckAuthorizationAccessToChannel(userID, channelID string) (bool, error) {
	channel, err := as.channelRepository.Find(channelID)
	if err != nil {
		return false, err
	}
	if channel == nil {
		return false, nil
	}

	userGuild, err := as.userGuildsRepository.Find(userID, *channel.GuildID)
	if err != nil {
		return false, err
	}
	if userGuild == nil {
		return false, nil
	}

	if !channel.Private {
		return true, nil
	}

	userChannel, err := as.userChannelsRepository.Find(userID, channelID)
	if err != nil {
		return false, err
	}
	if userChannel == nil {
		return false, nil
	}

	return true, nil
}
