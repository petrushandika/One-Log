package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashAPIKey creates a secure SHA-256 hash of a string value.
// We use this to hash API keys before saving them to the database.
func HashAPIKey(apiKey string) string {
	hasher := sha256.New()
	hasher.Write([]byte(apiKey))
	return hex.EncodeToString(hasher.Sum(nil))
}
