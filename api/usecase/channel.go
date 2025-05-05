package usecase

import (
	"chat_back/domain/model"
	"chat_back/domain/repository"

	"gorm.io/gorm"
)

type ChannelUsecase interface {
	Insert(db *gorm.DB, name string, description string, public bool, ownerID string, guildID *string) (*model.Channel, error)
	GetByID(db *gorm.DB, id string) (*model.Channel, error)
	GetAllInGuild(db *gorm.DB, guildID *string) ([]model.Channel, error)
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

func (cu channelUseCase) Insert(db *gorm.DB, name string, description string, public bool, ownerID string, guildID *string) (*model.Channel, error) {
	channel, err := cu.channelRepository.Insert(db, name, description, public, guildID)
	if err != nil {
		return nil, err
	}

	userChannels, err := cu.userChannelsRespository.Insert(db, ownerID, channel.ID)
	if err != nil {
		return nil, err
	}
	if userChannels == nil {
		return nil, err
	}
	return channel, nil
}

func (cu channelUseCase) GetByID(db *gorm.DB, id string) (*model.Channel, error) {
	channel, err := cu.channelRepository.GetByID(db, id)
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (cu channelUseCase) GetAllInGuild(db *gorm.DB, guildID *string) ([]model.Channel, error) {
	channels, err := cu.channelRepository.GetAllInGuild(db, guildID)
	if err != nil {
		return nil, err
	}
	return channels, nil
}
