package cover_service

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/models"
)

func (s *Service) Create(tx *sqlx.Tx, cover models.Cover) (createdCover models.Cover, err error) {
	coverId, err := s.CoverRepo.Create(tx, cover)
	if err != nil {
		return models.Cover{}, err
	}

	createdCover, err = s.CoverRepo.Read(tx, coverId)
	if err != nil {
		return models.Cover{}, err
	}

	return createdCover, nil
}
