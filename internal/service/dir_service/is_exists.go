package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s Service) IsExists(tx *sqlx.Tx, dirID int) (exists bool, err error) {
	log.Debug().Int("dirId", dirID).Msg("Checking directory existence")

	exists, err = s.dirRepo.IsExists(tx, dirID)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirID).Msg("Failed to check directory existence")
		return false, err
	}

	log.Debug().Int("dirId", dirID).Bool("exists", exists).Msg("Directory existence checked")
	return exists, nil
}
