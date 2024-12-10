package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

func GeneratePathHash(url string) string {
	hasher := sha256.New()
	hasher.Write([]byte(url))
	return hex.EncodeToString(hasher.Sum(nil))[:32]
}
