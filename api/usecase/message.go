package usecase

import (
	"chat_back/domain/model"
	"chat_back/domain/repository"

	"gorm.io/gorm"
)

type MessageUsecase interface {
	Insert(db *gorm.DB, channelID string, userID string, content string) error
	GetByID(db *gorm.DB, ID string) (model.Message, error)
	GetAllInChannel(db *gorm.DB, channelID string) ([]model.Message, error)
}

type messageUseCase struct {
	messageRepository repository.MessageRepository
}

func NewMessageUsecase(messageRepository repository.MessageRepository) MessageUsecase {
	return &messageUseCase{
		messageRepository: messageRepository,
	}
}

func (mu messageUseCase) Insert(db *gorm.DB, channelID string, userID string, content string) error {
	err := mu.messageRepository.Insert(db, channelID, userID, content)
	if err != nil {
		return err
	}
	return nil
}

func (mu messageUseCase) GetByID(db *gorm.DB, ID string) (model.Message, error) {
	message, err := mu.messageRepository.GetByID(db, ID)
	if err != nil {
		return model.Message{}, err
	}
	return message, nil
}
func (mu messageUseCase) GetAllInChannel(db *gorm.DB, channelID string) ([]model.Message, error) {
	messages, err := mu.messageRepository.GetAllInChannel(db, channelID)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
