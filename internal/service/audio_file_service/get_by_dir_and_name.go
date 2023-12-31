package audio_file_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/model"
)

func (s *Service) GetByDirAndName(tx *sqlx.Tx, dirId int, name string) (audioFile model.AudioFile, err error) {
	log.Debug().Int("dirId", dirId).Str("name", name).Msg("Getting audio file")

	exists, err := s.AudioFileRepo.IsExistsByDirAndName(tx, dirId, name)
	if err != nil {
		log.Error().Int("dirId", dirId).Str("name", name).Msg("Failed to check audio file existence")
		return model.AudioFile{}, err
	}
	if !exists {
		log.Error().Int("dirId", dirId).Str("name", name).Msg("Audio file not found")
		return model.AudioFile{}, errors.NotFound{Resource: fmt.Sprintf("audioFile with dirId=%d and name=%s in database", dirId, name)}
	}

	audioFile, err = s.AudioFileRepo.ReadByDirAndName(tx, dirId, name)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Str("name", name).Msg("Failed to fetch audio file")
		return model.AudioFile{}, err
	}

	log.Debug().Int("dirId", dirId).Str("name", name).Msg("Audio file got successfully")
	return audioFile, nil
}
