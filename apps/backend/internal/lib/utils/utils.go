// Package utils contains utility functions that our application needs
package utils

import (
	"crypto/rand"
	"math/big"
)

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// GenerateID creates a random ID of a specified length
func GenerateID(length int) string {
	id := make([]byte, length)
	alphabetLength := big.NewInt(int64(len(alphabet)))

	for i := range id {
		n, err := rand.Int(rand.Reader, alphabetLength)
		if err != nil {
			return ""
		}
		id[i] = alphabet[n.Int64()]
	}

	if string(id) == "" {
		panic("failed to create id")
	}

	return string(id)
}
