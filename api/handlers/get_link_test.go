package handlers_test

import (
	"bytes"
	"encoding/json"
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
	// First create a link to test with
	createBody, _ := json.Marshal(map[string]interface{}{
		"link": "https://example.com",
	})
	createReq, _ := http.NewRequest("POST", "/", bytes.NewReader(createBody))

	_, response := handlers.CreateHash(createReq)

	// Extract hash from response
	var createResponse map[string]string
	json.Unmarshal([]byte(response), &createResponse)
	hash := createResponse["hash"]

	// Now test getting the link
	req := &http.Request{
		URL: &url.URL{
			RawQuery: "hash=" + hash,
		},
	}

	status, body := handlers.GetLink(req)

	if status != http.StatusOK {
		t.Fatalf("Existing hash should return a 200 status. Instead returned %d", status)
	}

	if body == "" {
		t.Fatalf("No body returned.")
	}
}
