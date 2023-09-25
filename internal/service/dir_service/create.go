package dir_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
	"music-files/internal/utils"
)

func (s *Service) Create(tx *sqlx.Tx, dir models.Directory) (err error) {
	log.Debug().
		Interface("parentDirId", dir.ParentDirId).
		Str("name", dir.Name).
		Msg("Adding new directory")

	dir.Name = utils.SanitizePath(dir.Name)

	existsOnDisc, err := utils.DirectoryExists(dir.Name)
	if err != nil {
		log.Error().Err(err).
			Interface("parentDirId", dir.ParentDirId).
			Str("name", dir.Name)
		return err
	}
	if !existsOnDisc {
		err := fmt.Errorf("directory does not exist on disk")
		log.Error().Err(err)
		return err
	}

	existsInDb, err := s.DirRepo.IsExists(tx, dir.ParentDirId, dir.Name)
	if err != nil {
		log.Error().Err(err).
			Interface("parentDirId", dir.ParentDirId).
			Str("name", dir.Name)
		return err
	}
	if existsInDb {
		log.Info().
			Msg("Directory already added")
		return nil
	}

	_, err = s.DirRepo.Create(tx, dir)
	if err != nil {
		log.Error().Err(err).
			Interface("parentDirId", dir.ParentDirId).
			Str("name", dir.Name)
		return err
	}

	log.Debug().
		Interface("parentDirId", dir.ParentDirId).
		Str("name", dir.Name).
		Msg("Directory added successfully")
	return nil
}
