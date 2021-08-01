package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"url-shortener/db"
	"url-shortener/utils"
)


type body struct {
	Link string
}

func generateHash() string {
	// TODO: dynamically chose the minimum number of characters depending
	// on the number of links stored in the database?
	// 4 => ~7M
	// 5 => ~380M
	// 6 => ~19B
	return utils.RandomString(5)
}

func CreateHash(r *http.Request) (int, string) {
	decoder := json.NewDecoder(r.Body)
	var data body
	err := decoder.Decode(&data)
	if err != nil {
		return http.StatusBadRequest, ""
	}

	link := data.Link
	if link == "" {
		return http.StatusBadRequest, ""
	}

	// FIXME: check link is safe
	isSafe := true

	if !isSafe {
		return http.StatusForbidden, ""
	}

	hash := generateHash()

	log.Printf("Saving link '%s' with hash '%s'", link, hash)

	if err := db.InsertLink(hash, link); err != nil {
		log.Printf("Error executing insert-link.sql: %v\n", err)
		return http.StatusInternalServerError, ""
	}

	return http.StatusCreated, `{"hash":"` + hash + `"}`
}
