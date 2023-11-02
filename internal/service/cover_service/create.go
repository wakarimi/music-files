package cover_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model"
)

func (s *Service) Create(tx *sqlx.Tx, cover model.Cover) (createdCover model.Cover, err error) {
	log.Debug().Interface("cover", cover).Msg("Creating new cover")

	coverId, err := s.CoverRepo.Create(tx, cover)
	if err != nil {
		log.Error().Err(err).Interface("cover", cover).Msg("Failed to create new cover")
		return model.Cover{}, err
	}

	createdCover, err = s.CoverRepo.Read(tx, coverId)
	if err != nil {
		log.Error().Err(err).Interface("cover", cover).Msg("Failed to read created cover")
		return model.Cover{}, err
	}

	log.Debug().Interface("cover", cover).Msg("Cover created successfully")
	return createdCover, nil
}
