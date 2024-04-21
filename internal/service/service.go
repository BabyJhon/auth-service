package service

import (
	"context"

	"github.com/BabyJhon/auth-service/internal/repos"
)

type Auth interface {
	//auth methods
	RefreshTokens(ctx context.Context, accessToken, base64RefreshToken string) (string, string, error)
	CreateTokens(ctx context.Context, guid string) (string, string, error)
}

type Service struct {
	Auth
}

func NewService(repo *repos.Repository) *Service {
	return &Service{
		Auth: NewAuthService(repo),
	}
}
