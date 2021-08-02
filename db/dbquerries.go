package db

import (
	"context"
	"io/ioutil"
	"log"
)

var insertLinkSQL []byte
var selectLinkSQL []byte
var deleteOldLinksSQL []byte

func initQuerries() {
	var err error

	insertLinkSQL, err = ioutil.ReadFile("db/sql/insert-link.sql")
	if err != nil {
		log.Fatalf("Error while reading file 'db/sql/insert-link.sql': %v\n", err)
	}

	selectLinkSQL, err = ioutil.ReadFile("db/sql/select-link.sql")
	if err != nil {
		log.Fatalf("Error while reading file 'db/sql/select-link.sql': %v\n", err)
	}

	deleteOldLinksSQL, err = ioutil.ReadFile("db/sql/delete-links-old.sql")
	if err != nil {
		log.Fatalf("Error while reading file 'db/sql/delete-links-old.sql': %v\n", err)
	}
}

func InsertLink(hash string, link string) error {
	dbpool, err := getPool()
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return err
	}
	defer dbpool.Close()

	_, err = dbpool.Exec(context.Background(), string(insertLinkSQL), hash, link)

	return err
}

func SelectLink(hash string) (string, error) {
	dbpool, err := getPool()
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return "", err
	}
	defer dbpool.Close()

	var link string
	err = dbpool.QueryRow(context.Background(), string(selectLinkSQL), hash).Scan(&link)

	return link, err
}

func DeleteOldLinks() error {
	dbpool, err := getPool()
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return err
	}
	defer dbpool.Close()

	_, err = dbpool.Exec(context.Background(), string(deleteOldLinksSQL))

	return err
}
