package main

import (
	"net/http"

	"./database"
	"./handler"

	"github.com/gorilla/mux"
)

func main() {

	database.Connect()

	r := mux.NewRouter()
	r.HandleFunc("/file/upload", handler.ReceiveFile)
	r.HandleFunc("/{id}", handler.RenderFile)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.ListenAndServe(":8001", r)

}
