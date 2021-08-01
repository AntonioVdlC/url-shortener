package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"url-shortener/api"
	"url-shortener/db"
)

func init() {
	rand.Seed(time.Now().UnixNano())

	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	db.InitDB()
}

func main() {
	// API
	http.HandleFunc("/api/link", api.Link)

	// Static assets
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))),
	)

	// Port
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
