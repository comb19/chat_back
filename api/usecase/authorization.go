package usecase

import (
	"chat_back/domain/repository"
	"context"
	"log/slog"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
)

type AuthorizationUsecase interface {
	CheckPermission(channelID string, token string) (*clerk.User, error)
}

type authorizationUsecase struct {
	userChannelsRepository repository.UserChannelsRepository
}

func NewAuthorizationUsecase(userChannelsRepository repository.UserChannelsRepository) AuthorizationUsecase {
	return &authorizationUsecase{
		userChannelsRepository: userChannelsRepository,
	}
}

func (au *authorizationUsecase) CheckPermission(channelID string, token string) (*clerk.User, error) {
	slog.Info("CheckPermission")

	ctx := context.Background()
	claims, err := jwt.Verify(ctx, &jwt.VerifyParams{
		Token: token,
	})
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	user, err := user.Get(ctx, claims.Subject)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	userChannels, err := au.userChannelsRepository.Find(user.ID, channelID)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}
	if userChannels == nil {
		slog.ErrorContext(ctx, "a relation between user and channel is not found")
		return nil, nil
	}

	return user, nil
}
