package media

import (
	"context"
	"time"

	"../database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FileModel struct for holding th file metadata
type FileModel struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	Path string             `json:"path,omitempty" bson:"path"`
}

// GetFileDetails returns the file details
func GetFileDetails(fileID string) FileModel {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	collection := database.GetCollection("files")

	var result = FileModel{}
	collection.FindOne(ctx, bson.M{"_id": database.StringToObjectID(fileID)}).Decode(&result)
	return result
}
