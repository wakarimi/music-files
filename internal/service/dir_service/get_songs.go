package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model"
)

func (s *Service) GetAudioFiles(tx *sqlx.Tx, dirId int) (audioFiles []model.AudioFile, err error) {
	log.Debug().Msg("Getting audioFiles in directory")

	audioFiles, err = s.AudioFileService.GetAllByDir(tx, dirId)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Msg("Failed to get audioFiles in directory")
		return make([]model.AudioFile, 0), err
	}

	log.Debug().Int("countOfAudioFiles", len(audioFiles)).Msg("AudioFile in directory got successfully")
	return audioFiles, nil
}
