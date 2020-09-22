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
	r.Handle("/file/upload", RecoverWrap(http.HandlerFunc(handler.ReceiveFile))).Methods("POST")
	r.Handle("/{id}", RecoverWrap(http.HandlerFunc(handler.RenderFile))).Methods("GET")

	http.ListenAndServe(":8001", r)

}
