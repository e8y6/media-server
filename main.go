package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"./database"
	"./handler"

	"github.com/gorilla/mux"
)

func main() {

	database.Connect()

	r := mux.NewRouter()
	r.Handle("/file/upload", RecoverWrap(http.HandlerFunc(handler.ReceiveFile))).Methods("POST")
	r.Handle("/{ID}", RecoverWrap(http.HandlerFunc(handler.RenderFile))).Methods("GET")

	http.ListenAndServe(":8001", r)

}

func RecoverWrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}

				// TODO add sentry

				res := map[string]string{
					"error": err.Error(),
				}
				jsonResponse, _ := json.Marshal(res)

				w.Header().Set("Content-Type", "application/json")
				http.Error(w, string(jsonResponse), http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}
