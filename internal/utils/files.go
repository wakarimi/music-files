package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/h2non/filetype"
	"github.com/rs/zerolog/log"
	"os"
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

func IsMusicFile(absolutePath string) (isMusicFile bool, err error) {
	fileInfo, err := os.Stat(absolutePath)
	if err != nil {
		return false, err
	}
	if fileInfo.IsDir() {
		return false, err
	}

	file, err := os.Open(absolutePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	head := make([]byte, 261)
	file.Read(head)

	kind, _ := filetype.Match(head)
	if kind == filetype.Unknown {
		return false, nil
	}

	isMusicFile = kind.MIME.Value == "audio/mpeg" ||
		kind.MIME.Value == "audio/wav" ||
		kind.MIME.Value == "audio/flac" ||
		kind.MIME.Value == "audio/x-flac" ||
		kind.MIME.Value == "audio/aac" ||
		kind.MIME.Value == "audio/ogg" ||
		kind.MIME.Value == "audio/x-ms-wma" ||
		kind.MIME.Value == "audio/vnd.rn-realaudio" ||
		kind.MIME.Value == "audio/amr" ||
		kind.MIME.Value == "audio/mp4" ||
		kind.MIME.Value == "audio/alac" ||
		kind.MIME.Value == "audio/midi"

	return isMusicFile, nil
}

func IsImageFile(absolutePath string) (isImageFile bool, err error) {
	fileInfo, err := os.Stat(absolutePath)
	if err != nil {
		return false, err
	}
	if fileInfo.IsDir() {
		return false, err
	}

	file, err := os.Open(absolutePath)
	if err != nil {
		log.Warn().Err(err).Str("absolutePath", absolutePath).Msg("Failed to open file to check on image")
		return false, err
	}
	defer file.Close()

	head := make([]byte, 261)
	_, err = file.Read(head)
	if err != nil {
		log.Warn().Err(err).Str("absolutePath", absolutePath).Msg("Failed to read file to check on image")
		return false, err
	}

	kind, _ := filetype.Match(head)
	if kind == filetype.Unknown {
		return false, nil
	}

	isImageFile = kind.MIME.Value == "image/jpeg" ||
		kind.MIME.Value == "image/png" ||
		kind.MIME.Value == "image/gif"

	return isImageFile, nil
}
