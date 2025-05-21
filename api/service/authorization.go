package service

import "chat_back/domain/repository"

type AuthorizationService interface {
	CheckAuthorizationAccessToChannel(userID, channelID string) (bool, error)
}

type authorizationService struct {
	userChannelsRepository repository.UserChannelsRepository
}

func NewAuthorizationService(userChannelsRepository repository.UserChannelsRepository) AuthorizationService {
	return &authorizationService{
		userChannelsRepository: userChannelsRepository,
	}
}

func (as authorizationService) CheckAuthorizationAccessToChannel(userID, channelID string) (bool, error) {
	userChannel, err := as.userChannelsRepository.Find(userID, channelID)
	if err != nil {
		return false, err
	}
	if userChannel != nil {
		return false, nil
	}
	return true, nil
}
