package usecase

import (
	"chat_back/domain/model"
	"chat_back/domain/repository"

	"gorm.io/gorm"
)

type MessageUsecase interface {
	Insert(db *gorm.DB, channelID string, userID string, content string) (*model.Message, error)
	GetByID(db *gorm.DB, ID string) (*model.Message, error)
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

func (mu messageUseCase) Insert(db *gorm.DB, channelID string, userID string, content string) (*model.Message, error) {
	return mu.messageRepository.Insert(db, channelID, userID, content)
}

func (mu messageUseCase) GetByID(db *gorm.DB, ID string) (*model.Message, error) {
	return mu.messageRepository.GetByID(db, ID)
}
func (mu messageUseCase) GetAllInChannel(db *gorm.DB, channelID string) ([]model.Message, error) {
	return mu.messageRepository.GetAllInChannel(db, channelID)
}
