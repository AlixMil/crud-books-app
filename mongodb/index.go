package mongodb

import "go.mongodb.org/mongo-driver/mongo"

type MongoDB struct {
	booksCollection *mongo.Collection
	usersCollection *mongo.Collection
	filesCollection *mongo.Collection
	db              *mongo.Database
}

func New() *MongoDB {
	return &MongoDB{}
}
