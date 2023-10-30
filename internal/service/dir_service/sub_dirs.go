package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/models"
)

func (s *Service) SubDirs(tx *sqlx.Tx, dirId int) (roots []models.Directory, err error) {
	log.Debug().Int("dirId", dirId).Msg("Getting subdirectories")

	exists, err := s.DirRepo.IsExists(tx, dirId)
	if err != nil {
		log.Error().Int("dirId", dirId).Err(err).Msg("Failed to check directory existence")
		return make([]models.Directory, 0), err
	}
	if !exists {
		log.Error().Int("dirId", dirId).Msg("Directory not found")
		return make([]models.Directory, 0), errors.NotFound{Resource: "directory in database"}
	}

	roots, err = s.DirRepo.ReadSubDirs(tx, dirId)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Msg("Failed to get subdirectories")
		return make([]models.Directory, 0), err
	}

	log.Debug().Int("dirId", dirId).Int("subDirsCount", len(roots)).Msg("Subdirectories got successfully")
	return roots, nil
}
