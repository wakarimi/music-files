package audio_file_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
)

func (s *Service) Delete(tx *sqlx.Tx, audioFileId int) (err error) {
	log.Debug().Int("audioFileId", audioFileId).Msg("Deleting audio file")

	exists, err := s.AudioFileRepo.IsExists(tx, audioFileId)
	if err != nil {
		log.Error().Int("audioFileId", audioFileId).Msg("Failed to check audio file existence")
		return err
	}
	if !exists {
		log.Error().Int("audioFileId", audioFileId).Msg("Audio file not found")
		return errors.NotFound{Resource: fmt.Sprintf("audio_file with id=%d", audioFileId)}
	}

	err = s.AudioFileRepo.Delete(tx, audioFileId)
	if err != nil {
		log.Error().Err(err).Int("audioFileId", audioFileId).Msg("Failed to delete audio file")
		return err
	}

	log.Debug().Int("audioFileId", audioFileId).Msg("Audio file deleted successfully")
	return nil
}
