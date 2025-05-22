package usecase

import (
	"chat_back/domain/model"
	"chat_back/domain/repository"
	"chat_back/service"
	"fmt"
)

type ChannelUsecase interface {
	Insert(name string, description string, public bool, ownerID string, guildID *string) (*model.Channel, error)
	Delete(id string) error
	GetByID(id string) (*model.Channel, error)
	AddUserToChannel(id string, userIDs []string) (*model.Channel, error)
	GetMessagesOfChannel(id, userID string) (*[]*model.Message, error)
}

type channelUseCase struct {
	userChannelsRespository repository.UserChannelsRepository
	channelRepository       repository.ChannelRepository
	messageRepository       repository.MessageRepository
	authorizationService    service.AuthorizationService
}

func NewChannelUsecase(userChannelsRepository repository.UserChannelsRepository, channelRepository repository.ChannelRepository, messageRepository repository.MessageRepository, authorizationService service.AuthorizationService) ChannelUsecase {
	return &channelUseCase{
		userChannelsRespository: userChannelsRepository,
		channelRepository:       channelRepository,
		messageRepository:       messageRepository,
		authorizationService:    authorizationService,
	}
}

func (cu channelUseCase) Insert(name string, description string, public bool, ownerID string, guildID *string) (*model.Channel, error) {
	channel, err := cu.channelRepository.Insert(name, description, public, guildID)
	if err != nil {
		return nil, err
	}

	userChannels, err := cu.userChannelsRespository.Insert(ownerID, channel.ID)
	if err != nil {
		return nil, err
	}
	if userChannels == nil {
		return nil, err
	}
	return channel, nil
}

func (cu channelUseCase) Delete(id string) error {
	err := cu.channelRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (cu channelUseCase) GetByID(id string) (*model.Channel, error) {
	channel, err := cu.channelRepository.Find(id)
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (cu channelUseCase) AddUserToChannel(id string, userIDs []string) (*model.Channel, error) {
	channel, err := cu.channelRepository.Find(id)
	if err != nil {
		return nil, err
	}
	if channel == nil {
		return nil, fmt.Errorf("channel %s is not found", id)
	}

	for _, userID := range userIDs {
		userChannel, err := cu.userChannelsRespository.Insert(userID, id)
		if err != nil || userChannel == nil {
			continue
		}
	}
	return channel, nil
}

func (cu channelUseCase) GetMessagesOfChannel(id, userID string) (*[]*model.Message, error) {
	isAuthorized, err := cu.authorizationService.CheckAuthorizationAccessToChannel(userID, id)
	if err != nil {
		return nil, err
	}
	if !isAuthorized {
		return nil, nil
	}

	messages, err := cu.messageRepository.FindAllInChannel(id)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
