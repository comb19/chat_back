package usecase

import (
	"chat_back/domain/model"
	"chat_back/domain/repository"
)

type ChannelUsecase interface {
	Insert(name string, description string, public bool, ownerID string, guildID *string) (*model.Channel, error)
	GetByID(id string) (*model.Channel, error)
	GetAllInGuild(guildID *string) ([]model.Channel, error)
}

type channelUseCase struct {
	userChannelsRespository repository.UserChannelsRepository
	channelRepository       repository.ChannelRepository
}

func NewChannelUsecase(userChannelsRepository repository.UserChannelsRepository, channelRepository repository.ChannelRepository) ChannelUsecase {
	return &channelUseCase{
		userChannelsRespository: userChannelsRepository,
		channelRepository:       channelRepository,
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

func (cu channelUseCase) GetByID(id string) (*model.Channel, error) {
	channel, err := cu.channelRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (cu channelUseCase) GetAllInGuild(guildID *string) ([]model.Channel, error) {
	channels, err := cu.channelRepository.GetAllInGuild(guildID)
	if err != nil {
		return nil, err
	}
	return channels, nil
}
