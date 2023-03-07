package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	DB *mongo.Database
}

func New(login, pwd, dbName, host, port string) (*MongoDB, error) {
	ctx := context.Background()

	opts := options.Client()
	opts.ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s/", login, pwd, host, port))

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	return &MongoDB{DB: db}, nil
}
