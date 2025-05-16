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
}

type guildUseCase struct {
	guildRepository      repository.GuildRepository
	userGuildsRepository repository.UserGuildsRepository
	channelRepository    repository.ChannelRepository
}

func NewGuildUseCase(guildRepository repository.GuildRepository, userGuildsRepository repository.UserGuildsRepository, channelRepository repository.ChannelRepository) GuildUseCase {
	return &guildUseCase{
		guildRepository:      guildRepository,
		userGuildsRepository: userGuildsRepository,
		channelRepository:    channelRepository,
	}
}

func (gu guildUseCase) CreateNewGuild(name, description, ownerID string) (*model.Guild, error) {
	guild, err := gu.guildRepository.Insert(name, description)
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

	channels, err := gu.channelRepository.GetAllInGuild(&ID)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return channels, nil
}
