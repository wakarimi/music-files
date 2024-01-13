package audio_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/audio"
)

func (s Service) GetByDirAndName(tx *sqlx.Tx, dirID int, name string) (audioFile audio.Audio, err error) {
	log.Debug().Int("dirId", dirID).Str("name", name).Msg("Getting audio file")

	audioFile, err = s.audioRepo.ReadByDirAndName(tx, dirID, name)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirID).Str("name", name).Msg("Failed to fetch audio file")
		return audio.Audio{}, err
	}

	log.Debug().Int("dirId", dirID).Str("name", name).Msg("Audio file got successfully")
	return audioFile, nil
}
