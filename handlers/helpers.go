package handlers

import (
	"crypto/sha256"
	"encoding/hex"
)

func generateHash(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
