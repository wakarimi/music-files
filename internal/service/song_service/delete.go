package song_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"music-files/internal/errors"
)

func (s *Service) Delete(tx *sqlx.Tx, songId int) (err error) {
	exists, err := s.SongRepo.IsExists(tx, songId)
	if err != nil {
		return err
	}
	if !exists {
		return errors.NotFound{Resource: fmt.Sprintf("song with id=%d in database", songId)}
	}

	err = s.SongRepo.Delete(tx, songId)
	if err != nil {
		return err
	}

	return nil
}
