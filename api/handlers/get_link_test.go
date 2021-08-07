package handlers_test

import (
	"net/http"
	"net/url"
	"testing"

	"url-shortener/api/handlers"
)

func TestGetLinkNoHash(t *testing.T) {
	req := &http.Request{
		URL: &url.URL{
			RawQuery: "",
		},
	}

	status, _ := handlers.GetLink(req)

	if status != http.StatusBadRequest {
		t.Fatalf("No hash should return a 400 error. Instead returned %d", status)
	}
}

func TestGetLinkNotFound(t *testing.T) {
	req := &http.Request{
		URL: &url.URL{
			RawQuery: "hash=notfound",
		},
	}

	status, _ := handlers.GetLink(req)

	if status != http.StatusNotFound {
		t.Fatalf("No link should return a 404 error. Instead returned %d", status)
	}
}

func TestGetLink(t *testing.T) {
	req := &http.Request{
		URL: &url.URL{
			RawQuery: "hash=exists",
		},
	}

	status, body := handlers.GetLink(req)

	if status != http.StatusOK {
		t.Fatalf("No hash should return a 404 error. Instead returned %d", status)
	}

	if body == "" {
		t.Fatalf("No body returned.")
	}
}
