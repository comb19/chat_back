package usecase

import (
	"chat_back/domain/repository"
	"context"
	"fmt"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"gorm.io/gorm"
)

type AuthorizationUsecase interface {
	CheckPermission(db *gorm.DB, channelID string, token string) (*clerk.User, error)
}

type authorizationUsecase struct {
	userChannelsRepository repository.UserChannelsRepository
}

func NewAuthorizationUsecase(userChannelsRepository repository.UserChannelsRepository) AuthorizationUsecase {
	return &authorizationUsecase{
		userChannelsRepository: userChannelsRepository,
	}
}

func (au *authorizationUsecase) CheckPermission(db *gorm.DB, channelID string, token string) (*clerk.User, error) {
	fmt.Println("CheckPermission")
	context := context.Background()
	claims, err := jwt.Verify(context, &jwt.VerifyParams{
		Token: token,
	})
	if err != nil {
		return nil, err
	}
	user, err := user.Get(context, claims.Subject)
	if err != nil {
		return nil, err
	}

	userChannels, err := au.userChannelsRepository.Find(db, user.ID, channelID)
	if err != nil {
		return nil, err
	}
	if userChannels == nil {
		return nil, nil
	}
	fmt.Println("userChannels", userChannels)

	return user, nil
}
