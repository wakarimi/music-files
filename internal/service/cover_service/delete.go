package cover_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s Service) Delete(tx *sqlx.Tx, coverID int) error {
	log.Debug().Int("coverId", coverID).Msg("Deleting cover")

	err := s.coverRepo.Delete(tx, coverID)
	if err != nil {
		log.Error().Err(err).Int("coverId", coverID).Msg("Failed to delete cover")
		return err
	}

	log.Debug().Int("coverId", coverID).Msg("cover deleted")
	return nil
}
