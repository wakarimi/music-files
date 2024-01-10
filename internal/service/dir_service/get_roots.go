package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/directory"
)

func (s Service) GetRoots(tx *sqlx.Tx) (roots []directory.Directory, err error) {
	log.Debug().Msg("Getting root directories")

	roots, err = s.dirRepo.ReadRoots(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get root directories")
		return make([]directory.Directory, 0), err
	}

	log.Debug().Int("countOfRootDirs", len(roots)).Msg("Root directories got successfully")
	return roots, nil
}
