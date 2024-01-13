package audio_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/audio"
)

func (s Service) GetAllByDir(tx *sqlx.Tx, dirID int) (audios []audio.Audio, err error) {
	log.Debug().Int("dirId", dirID).Msg("Fetching audio")

	audios, err = s.audioRepo.ReadAllByDir(tx, dirID)
	if err != nil {
		log.Error().Int("dirId", dirID).Err(err).Msg("Failed to fetch all audio")
		return make([]audio.Audio, 0), err
	}

	log.Debug().Int("dirId", dirID).Int("countOfAudio", len(audios)).Msg("All audio files fetched successfully")
	return audios, nil
}
