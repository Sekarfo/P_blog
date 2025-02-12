package testutils

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

// TestPasswordHashing manually checks if stored password matches
func TestPasswordHashing(t *testing.T) {
	storedHash := "$2a$10$NDv9RczT3BKHIjMDLkXyc.obsdA8dQVugTWejFkI5CKH7iFrUCvZW" // Replace with actual stored hash
	inputPassword := "dias"                                                      // Replace with the actual password you used at registration

	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(inputPassword))
	if err != nil {
		t.Errorf("❌ Password does NOT match: %v", err)
	} else {
		fmt.Println("✅ Password matches!")
	}
}
