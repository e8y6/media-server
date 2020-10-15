package internalapi

import (
	"net/http"

	"../../media"

	"github.com/gorilla/mux"
)

func SaveFile(w http.ResponseWriter, r *http.Request) {

	urlParams := mux.Vars(r)
	fileID := urlParams["id"]
	fileObj := media.GetFileDetails(fileID)

	fileObj.Optimize()
	fileObj.MoveMediaSafe()

	fileObj.IsUsed = true

	fileObj.Save()

}
