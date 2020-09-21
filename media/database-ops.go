package media

import (
	"context"
	"time"

	"../database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetFileDetails returns the file details
func GetFileDetails(fileID string) FileModel {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	collection := database.GetCollection("files")

	var result = FileModel{}
	collection.FindOne(ctx, bson.M{"_id": database.StringToObjectID(fileID)}).Decode(&result)
	return result
}

func (fileObject *FileModel) Save() {

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	collection := database.GetCollection("files")

	_, err := collection.UpdateOne(
		ctx, bson.M{
			"_id": fileObject.ID,
		}, bson.M{"$set": fileObject}, options.Update().SetUpsert(true))

	if err != nil {
		panic(err)
	}

}