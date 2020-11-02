package main

import (
	"fmt"
	"net/http"

	"./config"
	"./database"
	"./handler/external"
	"./handler/internalapi"
	"./media/storage/cloudflare"
	"./misc/log"

	"github.com/gorilla/mux"
)

// Register initialization stuff here
func boot() {

	// Setup database connection
	database.Connect()

	// Get Signing keys from CF
	cloudflare.ObtainSigningKeys()

}

func main() {

	boot()

	externalRouter := mux.NewRouter()
	externalRouter.Handle("/{id}", http.HandlerFunc(external.ReceiveFileOptions)).Methods("OPTIONS")
	externalRouter.Handle("/file/upload", http.HandlerFunc(external.ReceiveFileOptions)).Methods("OPTIONS")

	externalRouter.Handle("/file/upload", Recover(http.HandlerFunc(external.ReceiveFile))).Methods("POST")
	externalRouter.Handle("/{id}", Recover(http.HandlerFunc(external.RenderFile))).Methods("GET")
	externalRouter.Handle("/file/{id}", Recover(http.HandlerFunc(external.RenderFile))).Methods("GET")
	externalRouter.Handle("/{id}/info", Recover(http.HandlerFunc(external.FileInfo))).Methods("GET")

	log.Info(fmt.Sprint("Public Server Starting on ", ":", config.APP_PORT_EXTERNAL))
	go http.ListenAndServe(fmt.Sprint(":", config.APP_PORT_EXTERNAL), externalRouter)

	internalRouter := mux.NewRouter()
	internalRouter.Handle("/file/{id}/save", Recover(http.HandlerFunc(internalapi.SaveFile))).Methods("POST")
	internalRouter.Handle("/file/{id}/delete", Recover(http.HandlerFunc(internalapi.DeleteFile))).Methods("POST")

	log.Info(fmt.Sprint("Internal Server Starting on ", ":", config.APP_PORT_INTERNAL))
	go http.ListenAndServe(fmt.Sprint(":", config.APP_PORT_INTERNAL), internalRouter)

	log.Info("================================================================================")

	select {}

}
