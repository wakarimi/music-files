package song_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"music-files/internal/errors"
	"music-files/internal/models"
)

func (s *Service) GetByDirAndName(tx *sqlx.Tx, dirId int, name string) (song models.Song, err error) {
	exists, err := s.SongRepo.IsExistsByDirAndName(tx, dirId, name)
	if err != nil {
		return models.Song{}, err
	}
	if !exists {
		return models.Song{}, errors.NotFound{Resource: fmt.Sprintf("song with dirId=%d and name=%s in database", dirId, name)}
	}

	song, err = s.SongRepo.ReadByDirAndName(tx, dirId, name)
	if err != nil {
		return models.Song{}, err
	}

	return song, nil
}
