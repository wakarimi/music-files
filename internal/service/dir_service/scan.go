package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/utils"
)

func (s *Service) Scan(tx *sqlx.Tx, dirId int) (err error) {
	log.Debug().Int("dirId", dirId).
		Msg("Scanning directory")

	absolutePath, err := s.AbsolutePath(tx, dirId)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	exists, err := utils.DirectoryExists(absolutePath)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	if !exists {
		// TODO: Удаление директории из базы данных
		log.Info().Err(err).Str("path", absolutePath).
			Msg("Directory not exists")
		return
	}

	log.Debug().Int("dirId", dirId).
		Msg("Directory scanned successfully")
	return nil
}
