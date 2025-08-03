package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"url-shortener/api"
	"url-shortener/db"
)

func TestMain(m *testing.M) {
	// Change to root directory so paths work correctly
	os.Chdir("..")
	
	// Load environment variables
	godotenv.Load(".env")
	
	// Initialize database
	db.Init()
	
	// Run tests
	code := m.Run()
	os.Exit(code)
}

// API Tests (testing through HTTP layer only)

func TestUnsupportedMethod(t *testing.T) {
	req, _ := http.NewRequest("PUT", "/api/link", nil)
	rr := httptest.NewRecorder()

	api.Link(rr, req)

	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("Unsupported method should return 501, got %v", status)
	}
}

func TestCreateLinkSuccess(t *testing.T) {
	body, _ := json.Marshal(map[string]interface{}{
		"link": "https://example.com",
	})
	req, _ := http.NewRequest("POST", "/api/link", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	api.Link(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("POST should return 201, got %v", status)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %v", contentType)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Response should be valid JSON: %v", err)
	}

	if hash, exists := response["hash"]; !exists || hash == "" {
		t.Errorf("Response should contain non-empty hash")
	}
}

func TestCreateAndRetrieveLink(t *testing.T) {
	// First create a link
	createBody, _ := json.Marshal(map[string]interface{}{
		"link": "https://example.com",
	})
	createReq, _ := http.NewRequest("POST", "/api/link", bytes.NewReader(createBody))
	createRR := httptest.NewRecorder()

	api.Link(createRR, createReq)

	var createResponse map[string]string
	json.Unmarshal(createRR.Body.Bytes(), &createResponse)
	hash := createResponse["hash"]

	// Now test getting the link
	getReq, _ := http.NewRequest("GET", "/api/link?hash="+hash, nil)
	getRR := httptest.NewRecorder()

	api.Link(getRR, getReq)

	if status := getRR.Code; status != http.StatusOK {
		t.Errorf("GET with valid hash should return 200, got %v", status)
	}

	contentType := getRR.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %v", contentType)
	}

	var getResponse map[string]string
	if err := json.Unmarshal(getRR.Body.Bytes(), &getResponse); err != nil {
		t.Errorf("Response should be valid JSON: %v", err)
	}

	if link, exists := getResponse["link"]; !exists || link != "https://example.com" {
		t.Errorf("Response should contain correct link, got %v", link)
	}
}

func TestGetLinkNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/link?hash=nonexistent", nil)
	rr := httptest.NewRecorder()

	api.Link(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("GET with invalid hash should return 404, got %v", status)
	}
}

func TestCreateLinkEmptyBody(t *testing.T) {
	body, _ := json.Marshal(map[string]interface{}{})
	req, _ := http.NewRequest("POST", "/api/link", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	api.Link(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Empty body should return 400, got %v", status)
	}
}

func TestCreateLinkEmptyLink(t *testing.T) {
	body, _ := json.Marshal(map[string]interface{}{
		"link": "",
	})
	req, _ := http.NewRequest("POST", "/api/link", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	api.Link(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Empty link should return 400, got %v", status)
	}
}

func TestCreateLinkInvalidURL(t *testing.T) {
	body, _ := json.Marshal(map[string]interface{}{
		"link": "invalid",
	})
	req, _ := http.NewRequest("POST", "/api/link", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	api.Link(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("Invalid URL should return 403, got %v", status)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %v", contentType)
	}

	if rr.Body.Len() == 0 {
		t.Errorf("Should return error message in body")
	}
}

func TestCreateLinkInvalidJSON(t *testing.T) {
	req, _ := http.NewRequest("POST", "/api/link", bytes.NewReader([]byte("invalid json")))
	rr := httptest.NewRecorder()

	api.Link(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Invalid JSON should return 400, got %v", status)
	}
}

func TestCreateLinkWrongFieldName(t *testing.T) {
	body, _ := json.Marshal(map[string]interface{}{
		"url": "https://example.com", // Wrong field name (should be "link")
	})
	req, _ := http.NewRequest("POST", "/api/link", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	api.Link(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Wrong field name should return 400, got %v", status)
	}
}

func TestCreateLinkDifferentURLTypes(t *testing.T) {
	testCases := []struct {
		name           string
		url            string
		expectedStatus int
	}{
		{"HTTPS with subdomain", "https://sub.example.com", http.StatusCreated},
		{"HTTPS with path", "https://example.com/path", http.StatusCreated},
		{"HTTPS with query", "https://example.com?query=1", http.StatusCreated},
		{"HTTP (insecure)", "http://example.com", http.StatusForbidden},
		{"FTP protocol", "ftp://example.com", http.StatusForbidden},
		{"No protocol", "example.com", http.StatusForbidden},
		{"No host", "https://", http.StatusForbidden},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(map[string]interface{}{
				"link": tc.url,
			})
			req, _ := http.NewRequest("POST", "/api/link", bytes.NewReader(body))
			rr := httptest.NewRecorder()

			api.Link(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("URL %s should return %d, got %d", tc.url, tc.expectedStatus, status)
			}
		})
	}
}

func TestCreateLinkResponseFormat(t *testing.T) {
	body, _ := json.Marshal(map[string]interface{}{
		"link": "https://example.com",
	})
	req, _ := http.NewRequest("POST", "/api/link", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	api.Link(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Should return 201 status, got %v", status)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Response should be valid JSON: %v", err)
	}

	hash, exists := response["hash"]
	if !exists {
		t.Errorf("Response should contain 'hash' field")
	}

	if len(hash) != 5 {
		t.Errorf("Hash should be 5 characters long, got %d", len(hash))
	}

	// Verify hash contains only valid characters
	validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for _, char := range hash {
		if !bytes.ContainsRune([]byte(validChars), char) {
			t.Errorf("Hash contains invalid character: %c", char)
		}
	}
}

func TestGetLinkNoHash(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/link", nil)
	rr := httptest.NewRecorder()

	api.Link(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("No hash should return 400, got %v", status)
	}
}

func TestGetLinkEmptyHash(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/link?hash=", nil)
	rr := httptest.NewRecorder()

	api.Link(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Empty hash should return 400, got %v", status)
	}
}

func TestGetLinkResponseFormat(t *testing.T) {
	// First create a link to test with
	createBody, _ := json.Marshal(map[string]interface{}{
		"link": "https://example.com/test",
	})
	createReq, _ := http.NewRequest("POST", "/api/link", bytes.NewReader(createBody))
	createRR := httptest.NewRecorder()

	api.Link(createRR, createReq)

	// Extract hash from response
	var createResponse map[string]string
	json.Unmarshal(createRR.Body.Bytes(), &createResponse)
	hash := createResponse["hash"]

	// Now test getting the link
	getReq, _ := http.NewRequest("GET", "/api/link?hash="+hash, nil)
	getRR := httptest.NewRecorder()

	api.Link(getRR, getReq)

	if status := getRR.Code; status != http.StatusOK {
		t.Errorf("Should return 200 status, got %v", status)
	}

	var getResponse map[string]string
	if err := json.Unmarshal(getRR.Body.Bytes(), &getResponse); err != nil {
		t.Errorf("Response should be valid JSON: %v", err)
	}

	link, exists := getResponse["link"]
	if !exists {
		t.Errorf("Response should contain 'link' field")
	}

	if link != "https://example.com/test" {
		t.Errorf("Response should contain correct link, got %s", link)
	}
}

func TestGetLinkDifferentHashFormats(t *testing.T) {
	testCases := []struct {
		name           string
		hash           string
		expectedStatus int
	}{
		{"Normal 5-char hash", "abcde", http.StatusNotFound},
		{"Long hash", "abcdefghijklmnop", http.StatusNotFound},
		{"Single char", "a", http.StatusNotFound},
		{"Special characters", "abc@#", http.StatusNotFound},
		{"Numbers", "12345", http.StatusNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/link?hash="+tc.hash, nil)
			rr := httptest.NewRecorder()

			api.Link(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("Hash %s should return %d, got %d", tc.hash, tc.expectedStatus, status)
			}
		})
	}
}