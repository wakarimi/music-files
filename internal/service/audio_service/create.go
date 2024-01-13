package audio_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/audio"
)

func (s Service) Create(tx *sqlx.Tx, audioToCreate audio.Audio) (int, error) {
	log.Debug().Interface("audioToCreate", audioToCreate).Msg("Creating new audio file")

	audioFileID, err := s.audioRepo.Create(tx, audioToCreate)
	if err != nil {
		log.Error().Err(err).Interface("audioToCreate", audioToCreate).Msg("Failed to create new audio file")
		return 0, err
	}

	log.Debug().Interface("audioToCreate", audioToCreate).Msg("Audio file created successfully")
	return audioFileID, nil
}
