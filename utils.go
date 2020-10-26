package main

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
	"strings"

	"./misc/exceptions"
	"./misc/log"
)

type ErrorResponse struct {
	Message string      `json:"message"`
	Trace   interface{} `json:"trace,omitempty"`
}

func GetResponse(err interface{}) (e ErrorResponse, statusCode int) {

	e.Message = "Some error ocurred"
	statusCode = 500

	switch t := err.(type) {
	case string:
		e.Message = t
	case error:
		e.Message = t.Error()
	case exceptions.Exception:
		e.Message = t.Message
		statusCode = t.GetStatusCode()
	}

	e.Trace = strings.Split(string(debug.Stack()), "\n\t")[4:]
	log.Error("Global exception handler caught an exception ", e.Message, e.Trace)

	return e, statusCode

}

// RecoverWrap will provide a gaceful death for panics
func Recover(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			r := recover()

			if r == nil {
				return
			}

			response, statusCode := GetResponse(r)
			jsonResponse, _ := json.Marshal(response)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			w.Write(jsonResponse)

		}()
		h.ServeHTTP(w, r)
	})
}
