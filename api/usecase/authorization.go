package usecase

import (
	"chat_back/service"
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
	authorizationService service.AuthorizationService
}

func NewAuthorizationUsecase(authorizationService service.AuthorizationService) AuthorizationUsecase {
	return &authorizationUsecase{
		authorizationService: authorizationService,
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

	isAuthorized, err := au.authorizationService.CheckAuthorizationAccessToChannel(user.ID, channelID)
	if err != nil {
		return nil, err
	}
	if !isAuthorized {
		return nil, nil
	}

	return user, nil
}
