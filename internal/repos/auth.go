package repos

import (
	"context"

	"github.com/BabyJhon/auth-service/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepo struct {
	mongoCollection *mongo.Collection
}

func NewAuthRepo(client *mongo.Client) *AuthRepo {
	return &AuthRepo{
		mongoCollection: client.Database("core").Collection("refreshSessions"),
	}
}

func (a *AuthRepo) AddSession(ctx context.Context, session entity.Session) error {
	_, err := a.mongoCollection.InsertOne(ctx, session)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthRepo) FindSessionsByGUID(ctx context.Context, guid string) ([]entity.Session, error) {
	cursor, err := a.mongoCollection.Find(ctx, bson.D{{"guid", guid}}) //вернет все сесии с таким гуидом
	if err != nil {
		return nil, err
	}

	var res []entity.Session
	if err = cursor.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (a *AuthRepo) DeleteSession(ctx context.Context, session entity.Session) error {
	_, err := a.mongoCollection.DeleteOne(ctx, bson.M{"refresh_token_hash": session.RefreshTokenHash})
	if err != nil {
		return err
	}
	return nil
}
