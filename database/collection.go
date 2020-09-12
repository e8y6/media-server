package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// DatabaseName Name of the database
	DatabaseName string = "johny_walker"
)

// GetCollection returns the collection
func GetCollection(collectionName string) *mongo.Collection {
	collection := client.Database(DatabaseName).Collection(collectionName)
	return collection
}
