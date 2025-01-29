package utils

import (
	"testing"
)

// go test ./utils -v

func TestHashPassword(t *testing.T) {
	password := "securepassword"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// is not empty
	if hashedPassword == "" {
		t.Error("Hashed password should not be empty")
	}
}

func TestComparePasswords(t *testing.T) {
	password := "securepassword"
	hashedPassword, _ := HashPassword(password)

	// Compare matching passwords
	if err := ComparePasswords(hashedPassword, password); err != nil {
		t.Errorf("Passwords should match, but got error: %v", err)
	}

	// Compare non-matching passwords
	if err := ComparePasswords(hashedPassword, "wrongpassword"); err == nil {
		t.Error("Passwords should not match, but no error returned")
	}
}
