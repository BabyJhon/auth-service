package repos

import (
	"context"
	"fmt"

	"github.com/BabyJhon/auth-service/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepo struct {
	//TODO: add mongo db
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
	//для проверкт работы
	fmt.Printf("repos: res len is: %d\n", len(res))
	for i := 0; i < len(res); i++ {
		fmt.Println(res[i].GUID)
	}
	return res, nil
}
