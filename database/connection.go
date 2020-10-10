package database

import (
	"context"
	"time"

	"../config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
)

// Connect to database
func Connect() {
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(config.DATABASE_CONNECTION_URI))
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
}

// GetConnection : Returns the database connection
func GetConnection() *mongo.Client {
	if client == nil {
		panic("Unable to connect..")
	}
	return client
}
