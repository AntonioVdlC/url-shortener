package api

import (
	"net/http"

	"url-shortener/api/handlers"
)

func Link(w http.ResponseWriter, r *http.Request) {
	var status int
	var data string

	switch r.Method {
	case "POST":
		status, data = handlers.CreateHash(r)
	case "GET":
		status, data = handlers.GetLink(r)
	default:
		status = http.StatusNotImplemented
		data = ""
	}

	w.WriteHeader(status)

	if data != "" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(data))
	}
}
