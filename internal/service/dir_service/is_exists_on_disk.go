package dir_service

import (
	"github.com/rs/zerolog/log"
	"os"
)

func (s Service) IsExistsOnDisk(path string) (bool, error) {
	log.Debug().Str("path", path).Msg("Checking directory on disk")

	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Error().Str("path", path).Msg("Directory not found on disk")
		return false, nil
	} else if err != nil {
		log.Error().Str("path", path).Msg("Failed to check directory existence")
		return false, err
	}

	if !fileInfo.IsDir() {
		log.Info().Str("path", path).Msg("The path to the file is specified instead of the directory")
		return false, nil
	}

	return true, nil
}
