package media

import (
	"context"
	"time"

	"../database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMedia(path string, fileName string, fileType string, userId string, privacy int) FileModel {

	return FileModel{
		Path:         path,
		UserID:       database.StringToObjectID(userId),
		FileType:     fileType,
		ID:           primitive.NewObjectID(),
		OriginalName: fileName,
		Privacy:      int8(privacy),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

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

func (fileObject *FileModel) Save() {

	fileObject.ID = primitive.NewObjectID()

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
