package cover_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"music-files/internal/errors"
)

func (s *Service) Delete(tx *sqlx.Tx, coverId int) (err error) {
	exists, err := s.CoverRepo.IsExists(tx, coverId)
	if err != nil {
		return err
	}
	if !exists {
		return errors.NotFound{Resource: fmt.Sprintf("cover with id=%d in database", coverId)}
	}

	err = s.CoverRepo.Delete(tx, coverId)
	if err != nil {
		return err
	}

	return nil
}
