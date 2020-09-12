package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"./database"
	"./media"

	"github.com/gorilla/mux"
)

func fileHandler(w http.ResponseWriter, r *http.Request) {

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

func main() {

	database.Connect()

	r := mux.NewRouter()
	r.HandleFunc("/{id}", fileHandler)

	http.ListenAndServe(":8001", r)

}
