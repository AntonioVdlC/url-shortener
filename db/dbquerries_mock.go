// +build mock

package db

import (
	"errors"
)

func initQuerries() {}

func InsertLink(hash string, link string) error {
	if link == "https://error.com" {
		return errors.New("Error")
	}
	return nil
}

func SelectLink(hash string) (string, error) {
	if hash == "notfound" {
		return "", errors.New("Not Found")
	}
	
	return "https://some-link", nil
}

func DeleteOldLinks() error {
	return nil
}
