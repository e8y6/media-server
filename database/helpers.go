package database

import (
	"../misc/exceptions"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StringToObjectID Convert string to mongo objectID
func StringToObjectID(objectIDString string) primitive.ObjectID {
	oid, error := primitive.ObjectIDFromHex(objectIDString)
	if error != nil {
		panic(exceptions.Exception{
			Message: "Unable to convert " + objectIDString + " usable form",
			Type:    exceptions.TYPE_PRECONDITION_FAILED,
		})
	}
	return oid
}
