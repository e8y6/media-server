package database

import "go.mongodb.org/mongo-driver/bson/primitive"

// StringToObjectID Convert string to mongo objectID
func StringToObjectID(objectIDString string) primitive.ObjectID {
	oid, _ := primitive.ObjectIDFromHex(objectIDString)
	return oid
}
