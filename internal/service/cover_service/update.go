package cover_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/models"
)

func (s *Service) Update(tx *sqlx.Tx, coverId int, cover models.Cover) (updatedCover models.Cover, err error) {
	log.Debug().Int("coverId", coverId).Msg("Updating cover")
	exists, err := s.CoverRepo.IsExists(tx, coverId)
	if err != nil {
		log.Debug().Int("coverId", coverId).Msg("Failed to check cover existence")
		return models.Cover{}, err
	}
	if !exists {
		log.Error().Int("coverId", coverId).Msg("Cover not found")
		return models.Cover{}, errors.NotFound{Resource: fmt.Sprintf("cover with coverId=%d in database", coverId)}
	}

	err = s.CoverRepo.Update(tx, coverId, cover)
	if err != nil {
		log.Debug().Int("coverId", coverId).Msg("Failed to update cover")
		return models.Cover{}, err
	}

	updatedCover, err = s.CoverRepo.Read(tx, coverId)
	if err != nil {
		log.Debug().Int("coverId", coverId).Msg("Failed to read updated cover")
		return models.Cover{}, err
	}

	log.Debug().Int("coverId", coverId).Msg("Cover updated successfully")
	return updatedCover, nil
}
