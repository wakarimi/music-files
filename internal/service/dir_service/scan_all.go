package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s *Service) ScanAll(tx *sqlx.Tx) (err error) {
	log.Debug().Msg("Scanning all root directories")

	roots, err := s.DirRepo.ReadRoots(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get root directories")
		return err
	}

	for _, root := range roots {
		err = s.Scan(tx, root.DirId)
		if err != nil {
			log.Error().Err(err).Int("rootDirId", root.DirId).Msg("Failed to scan directory")
			return err
		}
	}

	log.Debug().Msg("All root directories scanned successfully")
	return nil
}
