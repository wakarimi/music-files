package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (s *Service) GetRoots(tx *sqlx.Tx) (roots []models.Directory, err error) {
	log.Debug().Msg("Getting roots")

	roots, err = s.DirRepo.ReadRoots(tx)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to read root directories")
		return make([]models.Directory, 0), err
	}

	log.Debug().Int("rootsCount", len(roots)).Msg("Roots got successfully")
	return roots, nil
}
