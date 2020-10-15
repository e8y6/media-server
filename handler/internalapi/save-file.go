package internalapi

import (
	"fmt"
	"net/http"

	"../../media"

	"github.com/gorilla/mux"
)

func SaveFile(w http.ResponseWriter, r *http.Request) {

	urlParams := mux.Vars(r)
	fileID := urlParams["id"]
	fileObj := media.GetFileDetails(fileID)

	if fileObj.IsUsed {
		fmt.Println("File already moved to a safe place.")
		w.WriteHeader(201)
		return
	}

	fmt.Println(fileObj)

	fileObj.Optimize()
	fileObj.MoveMediaSafe()

	fmt.Println("File Moved to some safe place.")

	fileObj.IsUsed = true
	fileObj.Save()

}
