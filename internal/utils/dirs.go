package utils

import (
	"github.com/rs/zerolog/log"
	"os"
)

func IsDirectoryExistsOnDisk(path string) (directoryExists bool, err error) {
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		log.Info().Err(err).Str("path", path).Msg("Directory on disk not found")
		return false, nil
	} else if err != nil {
		log.Warn().Err(err).Str("path", path).Msg("Failed to check directory on disk")
		return false, err
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Warn().Err(err).Str("path", path).Msg("Failed to get object info")
		return false, err
	}
	if !fileInfo.IsDir() {
		log.Info().Err(err).Str("filepath", path).Msg("Trying to add a file instead of a folder")
		return false, err
	}

	return true, nil
}
