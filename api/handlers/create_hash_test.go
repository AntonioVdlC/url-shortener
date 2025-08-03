package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"url-shortener/api/handlers"
)


func TestCreateHashEmptyBody(t *testing.T) {
	body, _ := json.Marshal(map[string]interface{}{})
	req, _ := http.NewRequest("POST", "/", bytes.NewReader(body))

	status, _ := handlers.CreateHash(req)

	if status != http.StatusBadRequest {
		t.Fatalf("No link should return a 400 error. Instead returned %d", status)
	}
}

func TestCreateHashNoLink(t *testing.T) {
	body, _ := json.Marshal(map[string]interface{}{
		"link": "",
	})
	req, _ := http.NewRequest("POST", "/", bytes.NewReader(body))

	status, _ := handlers.CreateHash(req)

	if status != http.StatusBadRequest {
		t.Fatalf("No link should return a 400 error. Instead returned %d", status)
	}
}

func TestCreateHashInvalidLink(t *testing.T) {
	body, _ := json.Marshal(map[string]interface{}{
		"link": "invalid",
	})
	req, _ := http.NewRequest("POST", "/", bytes.NewReader(body))

	status, message := handlers.CreateHash(req)

	if status != http.StatusForbidden {
		t.Fatalf("Invalid link should return a 403 error. Instead returned %d", status)
	}
	if message == "" {
		t.Fatalf("Missing error message.")
	}
}

func TestCreateHashInvalidURL(t *testing.T) {
	body, _ := json.Marshal(map[string]interface{}{
		"link": "not-a-link",
	})
	req, _ := http.NewRequest("POST", "/", bytes.NewReader(body))

	status, _ := handlers.CreateHash(req)

	if status != http.StatusForbidden {
		t.Fatalf("Invalid URL should return a 403 error. Instead returned %d", status)
	}
}

func TestCreateHashSuccess(t *testing.T) {
	body, _ := json.Marshal(map[string]interface{}{
		"link": "https://www.google.com",
	})
	req, _ := http.NewRequest("POST", "/", bytes.NewReader(body))

	status, message := handlers.CreateHash(req)

	if status != http.StatusCreated {
		t.Fatalf("Insert link should return a 201 status code. Instead returned %d", status)
	}
	if message == "" {
		t.Fatalf("Missing return message.")
	}
}
