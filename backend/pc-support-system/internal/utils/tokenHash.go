package utils

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
)

func HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func VerifyRefreshToken(token, storedHash string) bool {
	hash := sha256.Sum256([]byte(token))

	expected, err := hex.DecodeString(storedHash)
	if err != nil || len(expected) != sha256.Size {
		return false
	}

	return subtle.ConstantTimeCompare(hash[:], expected) == 1
}
