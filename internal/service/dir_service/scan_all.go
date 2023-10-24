package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s *Service) ScanAll(tx *sqlx.Tx) (err error) {
	log.Debug().Msg("Scanning all root directories")

	roots, err := s.DirRepo.ReadRoots(tx)
	if err != nil {
		return err
	}

	for _, root := range roots {
		err := s.Scan(tx, root.DirId)
		if err != nil {
			return err
		}
	}

	log.Debug().Msg("All root directories scanned successfully")
	return nil
}
