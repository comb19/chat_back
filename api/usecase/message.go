package usecase

import (
	"chat_back/domain/model"
	"chat_back/domain/repository"
)

type MessageUsecase interface {
	Insert(channelID string, userID string, content string) (*model.Message, error)
	GetByID(ID string) (*model.Message, error)
	GetAllInChannel(channelID string) ([]model.Message, error)
}

type messageUseCase struct {
	messageRepository repository.MessageRepository
}

func NewMessageUsecase(messageRepository repository.MessageRepository) MessageUsecase {
	return &messageUseCase{
		messageRepository: messageRepository,
	}
}

func (mu messageUseCase) Insert(channelID string, userID string, content string) (*model.Message, error) {
	return mu.messageRepository.Insert(channelID, userID, content)
}

func (mu messageUseCase) GetByID(ID string) (*model.Message, error) {
	return mu.messageRepository.GetByID(ID)
}
func (mu messageUseCase) GetAllInChannel(channelID string) ([]model.Message, error) {
	return mu.messageRepository.GetAllInChannel(channelID)
}
