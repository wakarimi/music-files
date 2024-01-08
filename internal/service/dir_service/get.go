package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/directory"
)

func (s Service) Get(tx *sqlx.Tx, dirID int) (directory.Directory, error) {
	log.Debug().Int("dirId", dirID).Msg("Getting directory")

	dir, err := s.dirRepo.Read(tx, dirID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get directory")
		return directory.Directory{}, err
	}

	log.Debug().Msg("Directory got successfully")
	return dir, nil
}
