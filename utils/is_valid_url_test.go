package utils_test

import (
	"testing"

	"url-shortener/utils"
)

func TestValidUrl(t *testing.T) {
	url := "https://www.google.com"
	is_valid, err_msg := utils.IsValidUrl(url)
	if is_valid != true {
		t.Fatalf("URL %s is not valid: %s", url, err_msg)
	}
}

func TestInvalidUrlFormat(t *testing.T) {
	url := "hello, world!"
	is_valid, _ := utils.IsValidUrl(url)
	if is_valid == true {
		t.Fatalf("URL %s is valid", url)
	}
}

func TestInvalidUrlMissingScheme(t *testing.T) {
	url := "127.0.0.1"
	is_valid, _ := utils.IsValidUrl(url)
	if is_valid == true {
		t.Fatalf("URL %s is valid", url)
	}
}

func TestInvalidUrlNotSecureScheme(t *testing.T) {
	url := "http://www.google.com"
	is_valid, _ := utils.IsValidUrl(url)
	if is_valid == true {
		t.Fatalf("URL %s is valid", url)
	}
}

func TestInvalidUrlMissingHost(t *testing.T) {
	url := "https://"
	is_valid, _ := utils.IsValidUrl(url)
	if is_valid == true {
		t.Fatalf("URL %s is valid", url)
	}
}
