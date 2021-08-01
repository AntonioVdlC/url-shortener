package db

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

func InitDB() {
	// Get a pool
	dbpool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	// Check needed tables are created
	var initDone bool
	initCheckSQL, err := ioutil.ReadFile("db/sql/init-check.sql")
	if err != nil {
		log.Fatalf("Error while reading file 'db/sql/init-check.sql': %v\n", err)
	}
	err = dbpool.QueryRow(context.Background(), string(initCheckSQL)).Scan(&initDone)
	if err != nil {
		log.Printf("Error while checking table 'links': %v\n", err)
		initDone = false
	}

	// Create needed tables
	if !initDone {
		log.Println("Initialising tables in database ...")
		initSQL, err := ioutil.ReadFile("db/sql/init.sql")
		if err != nil {
			log.Fatalf("Error while reading file 'db/sql/init.sql': %v\n", err)
		}

		_, err = dbpool.Exec(context.Background(), string(initSQL))
		if err != nil {
			log.Fatalf("Error executing init.sql: %v\n", err)
		}
	}

	// Read other SQL files to be used in API handlers
	initQuerries()
}

func getPool() (*pgxpool.Pool, error) {
	return pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
}
