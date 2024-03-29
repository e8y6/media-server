package external

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"../../config"
	"../../database"
	"../../media"
	"../../misc/exceptions"
	"../../misc/log"
	"../../misc/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func findContentType(file *os.File) string {

	file.Seek(0, 0)

	buffer := make([]byte, 512) // seems like http.DetectContentType only use first 512 bytes
	n, err := file.Read(buffer)

	if err != nil {
		panic(err)
	}

	return http.DetectContentType(buffer[:n])

}

// ReceiveFile receives file
func ReceiveFile(w http.ResponseWriter, r *http.Request) {

	log.Info("Upload Started")

	// Upload and save File
	httpFile, header, err := r.FormFile("file")
	if err != nil {
		panic(exceptions.Exception{
			Message: "File not found in the request",
			Type:    exceptions.TYPE_BAD_REQUEST,
			Error:   err,
		})
	}
	defer httpFile.Close()

	path := utils.GenerateFileName(header.Filename)
	localFile, err := os.OpenFile(config.LOCAL_FOLDER+path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(exceptions.Exception{
			Type:  exceptions.TYPE_INTERNAL_ERROR,
			Error: err,
		})
	}
	defer localFile.Close()

	size, err := io.Copy(localFile, httpFile)
	if err != nil {
		panic(exceptions.Exception{
			Type:  exceptions.TYPE_INTERNAL_ERROR,
			Error: err,
		})
	}
	log.Info("Upload File saved to " + path)

	// Create DBentries
	fileType := findContentType(localFile) // TODO directly from multipart.File
	privacy, _ := strconv.Atoi(string(r.Form.Get("privacy")))
	myMedia := media.FileModel{
		UserID:       database.StringToObjectID(r.Form.Get("_id_user")),
		FileType:     fileType,
		IsUsed:       false,
		ID:           primitive.NewObjectID(),
		Size:         size,
		OriginalName: header.Filename,
		Bucket:       media.BUCKET_LOCAL,
		Privacy:      int8(privacy),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		BucketMeta: map[string]string{
			"path": path,
		},
	}
	myMedia.Save()

	log.Info("File object has been created for saved file ", myMedia.ID)

	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(myMedia)
	w.Write([]byte(response))

}
