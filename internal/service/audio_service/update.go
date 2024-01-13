package audio_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/audio"
)

func (s Service) Update(tx *sqlx.Tx, audioID int, audioToUpdate audio.Audio) (err error) {
	log.Debug().Int("audioId", audioID).Msg("Updating audio file")

	err = s.audioRepo.Update(tx, audioID, audioToUpdate)
	if err != nil {
		log.Debug().Int("audioId", audioID).Msg("Failed to update audio file")
		return err
	}

	log.Debug().Int("audioId", audioID).Msg("Audio file updated successfully")
	return nil
}
