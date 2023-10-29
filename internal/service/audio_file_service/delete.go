package audio_file_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"music-files/internal/errors"
)

func (s *Service) Delete(tx *sqlx.Tx, audioFileId int) (err error) {
	exists, err := s.AudioFileRepo.IsExists(tx, audioFileId)
	if err != nil {
		return err
	}
	if !exists {
		return errors.NotFound{Resource: fmt.Sprintf("audioFile with id=%d in database", audioFileId)}
	}

	err = s.AudioFileRepo.Delete(tx, audioFileId)
	if err != nil {
		return err
	}

	return nil
}
