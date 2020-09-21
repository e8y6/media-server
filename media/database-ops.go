package media

import (
	"context"
	"fmt"
	"time"

	"../database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewMedia(path string, fileName string, fileType string, userId string, privacy int) FileModel {

	return FileModel{
		Path:         path,
		UserID:       database.StringToObjectID(userId),
		FileType:     fileType,
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

func (fileObject *FileModel) SaveToDatabase() {

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	collection := database.GetCollection("files")

	insRes, err := collection.InsertOne(ctx, fileObject)
	if err != nil {
		panic(err)
	}
	fmt.Println()

	fileObject.ID = insRes.InsertedID.(primitive.ObjectID)

}

func (fileObject *FileModel) UpdateDatabaseEntry() {

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	collection := database.GetCollection("files")

	_, err := collection.UpdateOne(
		ctx, bson.M{
			"_id": fileObject.ID,
		}, bson.M{"$set": fileObject})

	if err != nil {
		panic(err)
	}

}
