package audio_file_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model"
)

func (s *Service) SearchBySha256(tx *sqlx.Tx, sha256 string) (audioFiles []model.AudioFile, err error) {
	log.Debug().Str("sha256", sha256).Msg("Fetching audio files by sha256")

	audioFiles, err = s.AudioFileRepo.ReadAllBySha256(tx, sha256)
	if err != nil {
		log.Error().Err(err).Str("sha256", sha256).Msg("Failed to fetch audio files")
		return make([]model.AudioFile, 0), err
	}

	log.Debug().Str("sha256", sha256).Int("countOfAudioFiles", len(audioFiles)).Msg("Audio files fetched successfully")
	return audioFiles, nil
}
