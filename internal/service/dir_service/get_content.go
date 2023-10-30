package dir_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/models"
)

func (s *Service) GetDir(tx *sqlx.Tx, dirId int) (dir models.Directory, err error) {
	log.Debug().Int("dirId", dirId).Msg("Getting content of directory")

	exists, err := s.DirRepo.IsExists(tx, dirId)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Msg("Failed to check directory existence")
		return models.Directory{}, err
	}
	if !exists {
		log.Error().Int("dirId", dirId).Msg("Directory not found")
		return models.Directory{}, errors.NotFound{Resource: fmt.Sprintf("directory with id=%d", dirId)}
	}

	dir, err = s.DirRepo.Read(tx, dirId)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Msg("Failed to read directory")
		return models.Directory{}, err
	}

	log.Debug().Int("dirId", dirId).Msg("Content of directory got successfully")
	return dir, nil
}
