package handlers

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hitpads/reado_ap/internal/db"
	"github.com/hitpads/reado_ap/internal/models"
	"github.com/hitpads/reado_ap/internal/users"
	"github.com/joho/godotenv"
)

// go test ./handlers -v

func TestRegisterHandler(t *testing.T) {
	// database connection
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	db.ConnectDB()
	defer db.DB.Exec("TRUNCATE users RESTART IDENTITY CASCADE") // reset DB

	// request payload
	payload := `{"name": "Test User", "email": "testuser@example.com", "password": "password123"}`
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	// response recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(users.RegisterHandler)

	// Execute
	handler.ServeHTTP(rr, req)

	// response status code
	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
	}

	// response body
	expected := "Registration successful! Please check your email to verify your account."
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("Unexpected response body: got %v, want %v", rr.Body.String(), expected)
	}

	// verify user saved in the database
	var user models.User
	if err := db.DB.Where("email = ?", "testuser@example.com").First(&user).Error; err != nil {
		t.Fatalf("User not found in database: %v", err)
	}

	if user.Name != "Test User" {
		t.Errorf("Unexpected user name: got %s, want %s", user.Name, "Test User")
	}
}
