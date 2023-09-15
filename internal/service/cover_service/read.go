package cover_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (s *Service) Read(tx *sqlx.Tx, coverId int) (cover models.Cover, err error) {
	cover, err = s.CoverRepo.ReadTx(tx, coverId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read cover")
		return models.Cover{}, err
	}

	return cover, nil
}
