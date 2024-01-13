package audio_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s Service) Delete(tx *sqlx.Tx, audioID int) error {
	log.Debug().Int("audioId", audioID).Msg("Deleting audio")

	err := s.audioRepo.Delete(tx, audioID)
	if err != nil {
		log.Error().Err(err).Int("audioId", audioID).Msg("Failed to delete audio")
		return err
	}

	log.Debug().Int("audioId", audioID).Msg("Audio deleted")
	return nil
}
