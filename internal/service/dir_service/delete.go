package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s Service) Delete(tx *sqlx.Tx, dirID int) (err error) {
	log.Debug().Int("dirId", dirID).Msg("Deleting directory")

	err = s.dirRepo.Delete(tx, dirID)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirID).Msg("Failed to delete directory")
		return err
	}

	log.Debug().Int("dirId", dirID).Msg("Directory deleted")
	return nil
}
