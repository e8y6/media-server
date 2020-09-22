package main

import (
	"fmt"
	"net/http"

	"./config"
	"./database"
	"./handler"

	"github.com/gorilla/mux"
)

func main() {

	database.Connect()

	r := mux.NewRouter()
	r.Handle("/file/upload", RecoverWrap(http.HandlerFunc(handler.ReceiveFile))).Methods("POST")
	r.Handle("/{id}", RecoverWrap(http.HandlerFunc(handler.RenderFile))).Methods("GET")

	fmt.Println(fmt.Sprint("App Starting on ", ":", config.APP_PORT))
	http.ListenAndServe(fmt.Sprint(":", config.APP_PORT), r)

}
