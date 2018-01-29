package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

// Computes an SHA 256 bit given the key string.
func computeHash(key string) string {
	hash := sha256.New()
	hash.Write([]byte(key))
	hashStr := hex.EncodeToString(hash.Sum(nil))
	return hashStr
}
