package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/rs/zerolog/log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CalculateSha256(filePath string) (hash string, err error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	hashBytes := sha256.Sum256(data)
	hash = hex.EncodeToString(hashBytes[:])
	return hash, nil
}

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
	imageExtensions := []string{".jpeg", ".jpg", ".png", ".gif"}
	ext := strings.ToLower(filepath.Ext(path))
	for _, imageExt := range imageExtensions {
		if ext == imageExt {
			return true
		}
	}
	return false
}

func GetAudioCodec(filepath string) (codecName string, err error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "a:0", "-show_entries", "stream=codec_name", "-of", "default=noprint_wrappers=1:nokey=1", filepath)
	codecBytes, err := cmd.Output()
	if err != nil {
		log.Error().Err(err).Msg("Failed to determine audio codec")
		return "", err
	}
	codecName = string(codecBytes)
	return codecName, nil
}
