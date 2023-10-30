package cover_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (s *Service) Create(tx *sqlx.Tx, cover models.Cover) (createdCover models.Cover, err error) {
	log.Debug().Interface("cover", cover).Msg("Creating new cover")

	coverId, err := s.CoverRepo.Create(tx, cover)
	if err != nil {
		log.Error().Err(err).Interface("cover", cover).Msg("Failed to create new cover")
		return models.Cover{}, err
	}

	createdCover, err = s.CoverRepo.Read(tx, coverId)
	if err != nil {
		log.Error().Err(err).Interface("cover", cover).Msg("Failed to read created cover")
		return models.Cover{}, err
	}

	log.Debug().Interface("cover", cover).Msg("Cover created successfully")
	return createdCover, nil
}
