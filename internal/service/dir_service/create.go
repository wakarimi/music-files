package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/directory"
)

func (s Service) Create(tx *sqlx.Tx, dirToCreate directory.Directory) (int, error) {
	log.Debug().Msg("Creating directory")

	createdDirID, err := s.dirRepo.Create(tx, dirToCreate)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create directory")
		return 0, err
	}

	log.Debug().Msg("Directory created successfully")
	return createdDirID, nil
}
