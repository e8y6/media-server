package database

import (
	"context"
	"time"

	"../config"
	"../misc/log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
)

// Connect to database
func Connect() {
	log.Info("Connecting to database.")
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(config.DATABASE_CONNECTION_URI))
	if err != nil {
		log.Fatal("Unable to create DB client.", err)
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Some error ocurred while connecting to DB server", err)
		panic(err)
	}
	log.Info("Connecting to database step completed successfully.")
}

// GetConnection : Returns the database connection
func GetConnection() *mongo.Client {
	if client == nil {
		panic("Unable to connect..")
	}
	return client
}
