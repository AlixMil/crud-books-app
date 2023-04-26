package mongodb

import (
	"context"
	"crud-books/config"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	usersCollectionName = "users"
	booksCollectionName = "books"
	filesCollectionName = "files"
)

func (m *MongoDB) Connect(cfg *config.Config) error {
	credential := options.Credential{
		AuthMechanism: "SCRAM-SHA-1",
		AuthSource:    cfg.DatabaseName,
		Username:      cfg.DatabaseLogin,
		Password:      cfg.DatabasePwd,
	}
	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s/", cfg.DatabaseHost, cfg.DatabasePort))
	opts.SetAuth(credential)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return fmt.Errorf("connecting to db failed, error: %w", err)
	}

	db := client.Database(cfg.DatabaseName)
	m.booksCollection = db.Collection(booksCollectionName)
	m.usersCollection = db.Collection(usersCollectionName)
	m.filesCollection = db.Collection(filesCollectionName)
	m.db = db

	return nil
}

func (m MongoDB) Ping() error {
	ctxPing, cancelPing := context.WithTimeout(context.Background(), time.Second*15)
	defer cancelPing()

	err := m.db.Client().Ping(ctxPing, readpref.Primary())

	if err != nil {
		return fmt.Errorf("ping error: %w", err)
	}

	return nil
}
