package internalapi

import (
	"net/http"

	"../../media"
	"../../misc/log"

	"github.com/gorilla/mux"
)

func SaveFile(w http.ResponseWriter, r *http.Request) {

	urlParams := mux.Vars(r)
	fileID := urlParams["id"]
	fileObj := media.GetFileDetails(fileID)

	if fileObj.IsUsed {
		log.Warn("File already moved to a safe place.", fileObj)
		w.WriteHeader(201)
		return
	}

	fileObj.Optimize()
	fileObj.MoveMediaSafe()

	log.Info("File Moved to some safe place.", fileID)

	fileObj.IsUsed = true
	fileObj.Save()

}
