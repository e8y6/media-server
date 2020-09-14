package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"../media"
	"github.com/gorilla/mux"
)

func RenderFile(w http.ResponseWriter, r *http.Request) {

	urlParams := mux.Vars(r)
	fileID := urlParams["id"]
	result := media.GetFileDetails(fileID)

	data, err := ioutil.ReadFile("persist/" + result.Path)
	if err != nil {
		panic(err)
	}
	w.Write(data)

	fmt.Println(result)
}
