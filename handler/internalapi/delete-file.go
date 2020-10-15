package internalapi

import (
	"net/http"
	"time"

	"../../media"

	"github.com/gorilla/mux"
)

func DeleteFile(w http.ResponseWriter, r *http.Request) {

	urlParams := mux.Vars(r)
	fileID := urlParams["id"]
	result := media.GetFileDetails(fileID)

	// Delete media from thirdparty integrations

	// Set as deleted in DB
	result.DeletedAt = time.Now()
	result.Save()

}
