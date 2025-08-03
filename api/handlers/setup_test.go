package handlers_test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"url-shortener/db"
)

func TestMain(m *testing.M) {
	// Change to root directory so paths work correctly
	os.Chdir("../..")
	
	// Load environment variables
	godotenv.Load(".env")
	
	// Initialize database
	db.Init()
	
	// Run tests
	code := m.Run()
	os.Exit(code)
}