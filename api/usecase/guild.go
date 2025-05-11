package usecase

import (
	"chat_back/domain/model"
	"chat_back/domain/repository"
)

type GuildUseCase interface {
	CreateNewGuild(name, description, ownerID string) (*model.Guild, error)
	GetGuildByID(id string) (*model.Guild, error)
}

type guildUseCase struct {
	guildRepository repository.GuildRepository
}

func NewGuildUseCase(guildRepository repository.GuildRepository) GuildUseCase {
	return &guildUseCase{
		guildRepository: guildRepository,
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
