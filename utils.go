package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"runtime/debug"
	"strings"

	"./misc/exceptions"
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
				case exceptions.Exception:
					res := map[string]interface{}{
						"error": t.Message,
						"trace": strings.Split(string(debug.Stack()), "\n\t")[4:],
					}
					jsonResponse, _ := json.Marshal(res)
					w.Header().Add("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(jsonResponse)
					return
				default:
					err = errors.New("Unknown error")
				}

				log.Error("Global exception handler caught an exception ", err.Error())

				res := map[string]interface{}{
					"error": err.Error(),
					"trace": strings.Split(string(debug.Stack()), "\n\t")[4:],
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
