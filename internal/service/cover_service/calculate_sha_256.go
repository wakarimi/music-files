package cover_service

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
)

func (s Service) CalculateSHA256(path string) (hash string, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	hashBytes := sha256.Sum256(data)
	hash = hex.EncodeToString(hashBytes[:])
	return hash, nil
}
