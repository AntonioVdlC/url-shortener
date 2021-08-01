package handlers

import (
	"log"
	"net/http"

	"url-shortener/db"
)

func GetLink(r *http.Request) (int, string) {
	hash := r.URL.Query().Get("hash")
	if hash == "" {
		return http.StatusBadRequest, ""
	}

	log.Printf("Redirecting from hash: '%s'", hash)

	link, err := db.SelectLink(hash)
	if err != nil {
		log.Printf("Link not found for hash: '%s'. Error: %v\n", hash, err)
		return http.StatusNotFound, ""
	}

	return http.StatusOK, `{"link": "` + link + `"}`
}
