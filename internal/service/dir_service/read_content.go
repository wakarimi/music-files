package dir_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (s *Service) ReadContent(tx *sqlx.Tx, dirId int) (dirs []models.Directory, err error) {
	parentDirExists, err := s.DirRepo.IsExists(tx, dirId)
	if err != nil {
		log.Error().Err(err)
		return make([]models.Directory, 0), err
	}
	if !parentDirExists {
		log.Error().Err(err)
		err = fmt.Errorf("directory not found")
		return make([]models.Directory, 0), err
	}

	dirs, err = s.DirRepo.ReadSubDirs(tx, dirId)
	if err != nil {
		log.Error().Err(err)
		return make([]models.Directory, 0), err
	}
	return dirs, nil
}
