package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"strings"
)

func (s Service) IsTracked(tx *sqlx.Tx, path string) (bool, error) {
	log.Debug().Msg("Checking if a directory is being tracked")

	rootDirs, err := s.dirRepo.ReadRoots(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read roots")
		return false, err
	}

	for _, rootDir := range rootDirs {
		if strings.HasPrefix(path, rootDir.Name) {
			log.Debug().Msg("Directory is being tracked")
			return true, nil
		}
	}

	log.Debug().Msg("The directory is not tracked")
	return false, nil
}
