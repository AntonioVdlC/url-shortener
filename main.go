package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/api/generate", generateHandler)
	http.HandleFunc("/api/redirect", redirectHandler)

	http.HandleFunc("/", indexHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/index.html")
}

type generateBody struct {
	Link string
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body)
		var data generateBody
		err := decoder.Decode(&data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		link := data.Link
		if link == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		
		// TODO: Generate hash for link
		hash := "blabla"

		// TODO: Save hash and link in db
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"hash":"` + hash + `"}`))
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		hash := r.URL.Query().Get("hash")
		if hash == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("Redirecting from hash: %s", hash)

		// TODO: get link corresponding to hash
		link := "https://www.google.com"

		// FIXME: error handling

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"link": "`+ link +`"}`))
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}
