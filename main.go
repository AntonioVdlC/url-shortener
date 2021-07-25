package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbpool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	var initDone bool
	initCheckSQL, err := ioutil.ReadFile("db/init-check.sql")
	if err != nil {
		log.Fatalf("Error while reading file 'db/init-check.sql': %v\n", err)
	}
	err = dbpool.QueryRow(context.Background(), string(initCheckSQL)).Scan(&initDone)
	if err != nil {
		log.Printf("Error while checking table 'links': %v\n", err)
		initDone = false
	}

	if !initDone {
		log.Println("Initialising tables in database ...")
		initSQL, err := ioutil.ReadFile("db/init.sql")
		if err != nil {
			log.Fatalf("Error while reading file 'db/init.sql': %v\n", err)
		}

		_, err = dbpool.Exec(context.Background(), string(initSQL))
		if err != nil {
			log.Fatalf("Error executing init.sql: %v\n", err)
		}
	}

	dbpool.Close()

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

		// FIXME: check link is safe
		isSafe := true
		isSageNextCheckAt := time.Now().Add(24 * time.Hour)
		hash := GenerateHash()

		log.Printf("Saving link '%s' with hash '%s'", link, hash)

		dbpool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
		if err != nil {
			log.Printf("Unable to connect to database: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer dbpool.Close()

		insertLinkSQL, err := ioutil.ReadFile("db/insert-link.sql")
		if err != nil {
			log.Printf("Error while reading file 'db/insert-link.sql': %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = dbpool.Exec(context.Background(), string(insertLinkSQL), hash, link, isSafe, isSageNextCheckAt)
		if err != nil {
			log.Printf("Error executing insert-link.sql: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

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

		log.Printf("Redirecting from hash: '%s'", hash)

		dbpool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
		if err != nil {
			log.Printf("Unable to connect to database: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer dbpool.Close()

		selectLinkSQL, err := ioutil.ReadFile("db/select-link.sql")
		if err != nil {
			log.Printf("Error while reading file 'db/select-link.sql': %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var link string
		err = dbpool.QueryRow(context.Background(), string(selectLinkSQL), hash).Scan(&link)
		if err != nil {
			log.Printf("Link not found for hash: '%s'. Error: %v\n", hash, err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"link": "` + link + `"}`))
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

// https://stackoverflow.com/a/31832326
const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(b)
}

func GenerateHash() string {
	// TODO: dynamically chose the minimum number of characters depending
	// on the number of links stored in the database?
	// 4 => ~7M
	// 5 => ~380M
	// 6 => ~19B
	return RandomString(5)
}
