package external

import (
	"encoding/json"
	"net/http"

	"../../media"

	"github.com/gorilla/mux"
)

// FileInfo returns the file info
func FileInfo(w http.ResponseWriter, r *http.Request) {

	urlParams := mux.Vars(r)
	fileID := urlParams["id"]
	result := media.GetFileDetails(fileID)

	jsonData, _ := json.Marshal(result)
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonData)

}
