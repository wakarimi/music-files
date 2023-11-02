package cover_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model"
)

func (s *Service) GetAllByDir(tx *sqlx.Tx, dirId int) (covers []model.Cover, err error) {
	log.Debug().Int("dirId", dirId).Msg("Fetching  covers")

	covers, err = s.CoverRepo.ReadAllByDir(tx, dirId)
	if err != nil {
		log.Error().Int("dirId", dirId).Err(err).Msg("Failed to fetch all covers")
		return make([]model.Cover, 0), err
	}

	log.Debug().Int("dirId", dirId).Int("countOfCovers", len(covers)).Msg("All covers fetched successfully")
	return covers, nil
}
