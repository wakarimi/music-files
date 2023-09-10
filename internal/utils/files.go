package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
)

func IsMusicFile(path string) bool {
	musicExtensions := []string{".aac", ".flac", ".m4a", ".mp3", ".ogg", ".opus", ".wav", ".wma"}
	ext := strings.ToLower(filepath.Ext(path))
	for _, musicExt := range musicExtensions {
		if ext == musicExt {
			return true
		}
	}
	return false
}

func IsImageFile(path string) bool {
	imageExtensions := []string{".jpeg", ".jpg", ".png", ".gif", ".bmp", ".tif", ".tiff", ".webp", ".ico", ".svg"}
	ext := strings.ToLower(filepath.Ext(path))
	for _, imageExt := range imageExtensions {
		if ext == imageExt {
			return true
		}
	}
	return false
}

func CalculateFileHash(filePath string) (hash string, err error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	hashBytes := sha256.Sum256(data)
	hash = hex.EncodeToString(hashBytes[:])
	return hash, nil
}
