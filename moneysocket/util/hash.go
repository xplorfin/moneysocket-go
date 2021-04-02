package util

import "crypto/sha256"

// CreateSha256Hash creates a hash from data.
func CreateSha256Hash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// CreateDoubleSha256 creates a double sha256 hash from data.
func CreateDoubleSha256(data []byte) []byte {
	fixed := CreateSha256Hash(CreateSha256Hash(data))
	return fixed[:]
}
