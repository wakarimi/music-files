package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (s *Service) ReadAll(tx *sqlx.Tx) (dirs []models.Directory, err error) {
	dirs, err = s.DirRepo.ReadAllTx(tx)
	if err != nil {
		log.Info().Msg("Failed to get all directories")
		return make([]models.Directory, 0), err
	}

	return dirs, err
}
