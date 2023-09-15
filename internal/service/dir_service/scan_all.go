package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s *Service) ScanAll(tx *sqlx.Tx) (err error) {
	dirs, err := s.DirRepo.ReadAllTx(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get directories")
		return err
	}

	for _, dir := range dirs {
		log.Debug().Int("dirId", dir.DirId).Msg("Scanning directory")
		s.dirScan(tx, dir)
	}

	return nil
}
