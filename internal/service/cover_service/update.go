package cover_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"music-files/internal/errors"
	"music-files/internal/models"
)

func (s *Service) Update(tx *sqlx.Tx, coverId int, cover models.Cover) (updatedCover models.Cover, err error) {
	exists, err := s.CoverRepo.IsExists(tx, coverId)
	if err != nil {
		return models.Cover{}, err
	}
	if !exists {
		return models.Cover{}, errors.NotFound{Resource: fmt.Sprintf("cover with coverId=%d in database", coverId)}
	}

	err = s.CoverRepo.Update(tx, coverId, cover)
	if err != nil {
		return models.Cover{}, err
	}

	updatedCover, err = s.CoverRepo.Read(tx, coverId)
	if err != nil {
		return models.Cover{}, err
	}

	return updatedCover, nil
}
