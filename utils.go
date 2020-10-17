package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"./misc/log"
)

// RecoverWrap will provide a gaceful death for panics
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

				log.Error("Global exception handler caught an exception ", err.Error())

				res := map[string]string{
					"error": err.Error(),
				}
				jsonResponse, _ := json.Marshal(res)

				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonResponse)

			}
		}()
		h.ServeHTTP(w, r)
	})
}
