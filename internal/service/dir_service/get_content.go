package dir_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/models"
)

func (s *Service) GetDir(tx *sqlx.Tx, dirId int) (dir models.Directory, err error) {
	log.Debug().Int("dirId", dirId).Msg("Getting directory")

	exists, err := s.DirRepo.IsExists(tx, dirId)
	if err != nil {
		return models.Directory{}, err
	}
	if !exists {
		err = errors.NotFound{Resource: fmt.Sprintf("directory with id=%d", dirId)}
		return models.Directory{}, err
	}

	dir, err = s.DirRepo.Read(tx, dirId)
	if err != nil {
		return models.Directory{}, err
	}

	log.Debug().Int("dirId", dirId).Msg("Directory got successfully")
	return dir, nil
}
