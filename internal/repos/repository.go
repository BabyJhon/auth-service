package repos

import (
	"context"

	"github.com/BabyJhon/auth-service/internal/entity"
	"go.mongodb.org/mongo-driver/mongo"
)

type Auth interface {
	AddSession(ctx context.Context, session entity.Session) error
	FindSessionsByGUID(ctx context.Context, guid string) ([]entity.Session, error)
	DeleteSession(ctx context.Context, session entity.Session) error 
}

type Repository struct {
	Auth
}

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		Auth: NewAuthRepo(client),
	}
}
