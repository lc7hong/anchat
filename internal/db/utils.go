package db

import (
	"crypto/sha256"
)

// HashKey creates a blind index for a public key
func HashKey(pubkey []byte) []byte {
	hash := sha256.Sum256(pubkey)
	return hash[:]
}
