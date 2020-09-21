package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"../database"
	"../media"
	"../utils"
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

	// Upload and save File
	httpFile, header, err := r.FormFile("file")
	if err != nil {
		panic("File not found in the request")
	}
	defer httpFile.Close()
	path := "uploads/" + utils.GenerateFileName(header.Filename)
	localFile, err := os.OpenFile("./persist/"+path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer localFile.Close()
	io.Copy(localFile, httpFile)
	// Upload complete

	// Create DBentries
	fileType := findContentType(localFile) // TODO directly from multipart.File
	privacy, _ := strconv.Atoi(string(r.Form.Get("privacy")))
	myMedia := media.FileModel{
		UserID:       database.StringToObjectID(r.Form.Get("_id_user")),
		FileType:     fileType,
		ID:           primitive.NewObjectID(),
		OriginalName: header.Filename,
		Bucket:       media.LOCAL,
		Privacy:      int8(privacy),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	myMedia.BucketMeta = map[string]string{
		"path": path,
	}
	myMedia.Save()

	myMedia.ProcessMedia()
	myMedia.MoveMediaSafe()

	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(myMedia)
	w.Write([]byte(response))

}