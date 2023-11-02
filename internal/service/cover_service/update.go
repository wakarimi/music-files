package cover_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/model"
)

func (s *Service) Update(tx *sqlx.Tx, coverId int, cover model.Cover) (updatedCover model.Cover, err error) {
	log.Debug().Int("coverId", coverId).Msg("Updating cover")
	exists, err := s.CoverRepo.IsExists(tx, coverId)
	if err != nil {
		log.Debug().Int("coverId", coverId).Msg("Failed to check cover existence")
		return model.Cover{}, err
	}
	if !exists {
		log.Error().Int("coverId", coverId).Msg("Cover not found")
		return model.Cover{}, errors.NotFound{Resource: fmt.Sprintf("cover with coverId=%d in database", coverId)}
	}

	err = s.CoverRepo.Update(tx, coverId, cover)
	if err != nil {
		log.Debug().Int("coverId", coverId).Msg("Failed to update cover")
		return model.Cover{}, err
	}

	updatedCover, err = s.CoverRepo.Read(tx, coverId)
	if err != nil {
		log.Debug().Int("coverId", coverId).Msg("Failed to read updated cover")
		return model.Cover{}, err
	}

	log.Debug().Int("coverId", coverId).Msg("Cover updated successfully")
	return updatedCover, nil
}
