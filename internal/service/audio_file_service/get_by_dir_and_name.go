package audio_file_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"music-files/internal/errors"
	"music-files/internal/models"
)

func (s *Service) GetByDirAndName(tx *sqlx.Tx, dirId int, name string) (audioFile models.AudioFile, err error) {
	exists, err := s.AudioFileRepo.IsExistsByDirAndName(tx, dirId, name)
	if err != nil {
		return models.AudioFile{}, err
	}
	if !exists {
		return models.AudioFile{}, errors.NotFound{Resource: fmt.Sprintf("audioFile with dirId=%d and name=%s in database", dirId, name)}
	}

	audioFile, err = s.AudioFileRepo.ReadByDirAndName(tx, dirId, name)
	if err != nil {
		return models.AudioFile{}, err
	}

	return audioFile, nil
}
