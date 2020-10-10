package database

import (
	"../config"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetCollection returns the collection
func GetCollection(collectionName string) *mongo.Collection {
	collection := client.Database(config.DATABASE_NAME).Collection(collectionName)
	return collection
}
