package cover_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/models"
)

func (s *Service) GetCover(tx *sqlx.Tx, coverId int) (cover models.Cover, err error) {
	log.Debug().Int("coverId", coverId).Msg("Getting cover")

	exists, err := s.CoverRepo.IsExists(tx, coverId)
	if err != nil {
		log.Error().Int("coverId", coverId).Msg("Failed to check cover existence")
		return models.Cover{}, err
	}
	if !exists {
		log.Error().Int("coverId", coverId).Msg("Cover not found")
		return models.Cover{}, errors.NotFound{Resource: fmt.Sprintf("cover with id=%d", coverId)}
	}

	cover, err = s.CoverRepo.Read(tx, coverId)
	if err != nil {
		log.Error().Err(err).Int("coverId", coverId).Msg("Failed to fetch cover")
		return models.Cover{}, err
	}

	log.Debug().Int("coverId", coverId).Msg("Cover got successfully")
	return cover, nil
}
