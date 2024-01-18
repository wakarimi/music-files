package audio_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/audio"
)

func (s Service) Get(tx *sqlx.Tx, audioID int) (audioFile audio.Audio, err error) {
	log.Debug().Int("audioId", audioID).Msg("Getting audio file")

	audioFile, err = s.audioRepo.Read(tx, audioID)
	if err != nil {
		log.Error().Err(err).Int("audioId", audioID).Msg("Failed to fetch audio file")
		return audio.Audio{}, err
	}

	log.Debug().Int("audioId", audioID).Msg("Audio file got successfully")
	return audioFile, nil
}
