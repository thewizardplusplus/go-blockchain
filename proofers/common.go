package proofers

import (
	"crypto/sha256"
)

func makeHash(data string) []byte {
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}
