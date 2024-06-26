package mongodbgo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewMongoDB(ctx context.Context) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://db:27017")) //"mongodb://127.0.0.1:27017" - при запуске монги отдельно и go run ./cmd

	if err != nil {
		return nil, err
	}

	return client, nil
}
