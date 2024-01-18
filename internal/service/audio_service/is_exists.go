package audio_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s Service) IsExists(tx *sqlx.Tx, audioID int) (exists bool, err error) {
	log.Debug().Int("audioId", audioID).Msg("Checking audio file existence")

	exists, err = s.audioRepo.IsExists(tx, audioID)
	if err != nil {
		log.Debug().Int("audioId", audioID).Msg("Failed to check audio file existence")
		return false, err
	}

	log.Debug().Int("audioId", audioID).Bool("exists", exists).Msg("Audio file existence checked successfully")
	return exists, nil
}
