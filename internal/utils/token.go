package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomToken generates a random token of the given byte length
func GenerateRandomToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		panic("Failed to generate random token")
	}
	return hex.EncodeToString(bytes)
}
