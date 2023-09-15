package dir_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (s *Service) Create(tx *sqlx.Tx, dir models.Directory) (err error) {
	exists, err := s.DirRepo.IsExistsByPathTx(tx, dir.Path)
	if err != nil {
		log.Debug().Msg("Failed to check directory existence")
		return err
	}
	if exists {
		err := fmt.Errorf("directory alrady exists")
		log.Info().Msg("Directory already added")
		return err
	}

	_, err = s.DirRepo.CreateTx(tx, dir)
	if err != nil {
		log.Info().Msg("Failed to add directory")
		return err
	}

	return nil
}
