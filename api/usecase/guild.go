package usecase

import (
	"chat_back/domain/model"
	"chat_back/domain/repository"
	"log/slog"
)

type GuildUseCase interface {
	CreateNewGuild(name, description, ownerID string) (*model.Guild, error)
	GetGuildByID(id string) (*model.Guild, error)
	GetGuildsOfUser(userID string) ([]*model.Guild, error)
	GetChannelsOfGuild(ID, userID string) ([]*model.Channel, error)
	CreateNewChannelInGuild(guildID, userID, name, description string) (*model.Channel, error)
}

type guildUseCase struct {
	guildRepository        repository.GuildRepository
	userGuildsRepository   repository.UserGuildsRepository
	channelRepository      repository.ChannelRepository
	userChannelsRepository repository.UserChannelsRepository
}

func NewGuildUseCase(guildRepository repository.GuildRepository, userGuildsRepository repository.UserGuildsRepository, channelRepository repository.ChannelRepository, userChannelsRepository repository.UserChannelsRepository) GuildUseCase {
	return &guildUseCase{
		guildRepository:        guildRepository,
		userGuildsRepository:   userGuildsRepository,
		channelRepository:      channelRepository,
		userChannelsRepository: userChannelsRepository,
	}
}

func (gu guildUseCase) CreateNewGuild(name, description, ownerID string) (*model.Guild, error) {
	guild, err := gu.guildRepository.Insert(name, description)
	if err != nil {
		return nil, err
	}
	_, err = gu.userGuildsRepository.Insert(ownerID, guild.ID)
	if err != nil {
		return nil, err
	}
	return guild, nil
}

func (gu guildUseCase) GetGuildByID(id string) (*model.Guild, error) {
	guild, err := gu.guildRepository.Find(id)
	if err != nil {
		return nil, err
	}
	return guild, nil
}

func (gu guildUseCase) GetGuildsOfUser(userID string) ([]*model.Guild, error) {
	return gu.guildRepository.FindOfUser(userID)
}

func (gu guildUseCase) GetChannelsOfGuild(ID, userID string) ([]*model.Channel, error) {
	userGuilds, err := gu.userGuildsRepository.Find(userID, ID)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	if userGuilds == nil {
		return nil, nil
	}

	channels, err := gu.channelRepository.FindAllInGuild(&ID)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return channels, nil
}

func (gu guildUseCase) CreateNewChannelInGuild(guildID, userID, name, description string) (*model.Channel, error) {
	channel, err := gu.channelRepository.Insert(name, description, false, &guildID)
	if err != nil {
		return nil, err
	}

	_, err = gu.userChannelsRepository.Insert(userID, channel.ID)
	if err != nil {
		return nil, err
	}

	return channel, nil
}
