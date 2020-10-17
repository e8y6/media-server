package external

import (
	"encoding/json"
	"net/http"
	"time"

	"../../media"

	"github.com/gorilla/mux"
)

const MODULE = "EXTERNAL_HANDLER"

type FileInfoResponse struct {
	FileType  string    `json:"type"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// FileInfo returns the file info
func FileInfo(w http.ResponseWriter, r *http.Request) {

	urlParams := mux.Vars(r)
	fileID := urlParams["id"]
	result := media.GetFileDetails(fileID)

	outInfo := FileInfoResponse{
		FileType:  result.FileType,
		Name:      result.OriginalName,
		CreatedAt: result.CreatedAt,
	}

	jsonData, _ := json.Marshal(outInfo)

	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonData)

}
