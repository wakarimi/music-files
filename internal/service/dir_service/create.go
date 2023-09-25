package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (s *Service) Create(tx *sqlx.Tx, dir models.Directory) (err error) {
	log.Debug().
		Interface("parentDirId", dir.ParentDirId).
		Str("name", dir.Name).
		Msg("Adding new directory")

	exists, err := s.DirRepo.IsExists(tx, dir.DirId, dir.Name)
	if err != nil {
		log.Error().Err(err).
			Interface("parentDirId", dir.ParentDirId).
			Str("name", dir.Name)
		return err
	}
	if exists {
		log.Info().
			Msg("Directory already added")
		return nil
	}

	// TODO: Проверить не заканчивается ли путь на /

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
