package cover_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
)

func (s *Service) Delete(tx *sqlx.Tx, coverId int) (err error) {
	log.Debug().Int("coverId", coverId).Msg("Deleting cover")

	exists, err := s.CoverRepo.IsExists(tx, coverId)
	if err != nil {
		log.Error().Int("coverId", coverId).Msg("Failed to check cover existence")
		return err
	}
	if !exists {
		log.Error().Int("coverId", coverId).Msg("Cover not found")
		return errors.NotFound{Resource: fmt.Sprintf("cover with id=%d", coverId)}
	}

	err = s.CoverRepo.Delete(tx, coverId)
	if err != nil {
		log.Error().Err(err).Int("coverId", coverId).Msg("Failed to delete cover")
		return err
	}

	log.Debug().Int("coverId", coverId).Msg("Cover deleted successfully")
	return nil
}
