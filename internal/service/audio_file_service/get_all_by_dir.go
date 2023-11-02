package audio_file_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model"
)

func (s *Service) GetAllByDir(tx *sqlx.Tx, dirId int) (audioFiles []model.AudioFile, err error) {
	log.Debug().Int("dirId", dirId).Msg("Fetching  audio files")

	audioFiles, err = s.AudioFileRepo.ReadAllByDir(tx, dirId)
	if err != nil {
		log.Error().Int("dirId", dirId).Err(err).Msg("Failed to fetch all audio files")
		return make([]model.AudioFile, 0), err
	}

	log.Debug().Int("dirId", dirId).Int("countOfAudioFiles", len(audioFiles)).Msg("All audio files fetched successfully")
	return audioFiles, nil
}
