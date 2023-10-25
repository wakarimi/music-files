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
		return make([]models.Directory, 0), err
	}
	if !exists {
		return make([]models.Directory, 0), errors.NotFound{Resource: "directory in database"}
	}

	roots, err = s.DirRepo.ReadSubDirs(tx, dirId)
	if err != nil {
		log.Warn().Int("dirId", dirId).Err(err).Msg("Failed to get subdirectories")
		return make([]models.Directory, 0), err
	}

	log.Debug().Int("dirId", dirId).Int("subDirsCount", len(roots)).Msg("Subdirectories got successfully")
	return roots, nil
}
