package util

import "crypto/sha256"

func CreateSha256Hash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func CreateDoubleSha256(data []byte) []byte {
	fixed := CreateSha256Hash(CreateSha256Hash(data))
	return fixed[:]
}
