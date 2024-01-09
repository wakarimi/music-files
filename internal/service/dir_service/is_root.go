package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s Service) IsRoot(tx *sqlx.Tx, dirID int) (isRoot bool, err error) {
	log.Debug().Int("dirId", dirID).Msg("Checking if the directory is the root")

	isRoot, err = s.dirRepo.IsRoot(tx, dirID)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirID).Msg("Could not check whether the directory is the root")
		return false, err
	}

	log.Debug().Int("dirId", dirID).Bool("isRoot", isRoot).Msg("Checking whether the directory is the root is complete")
	return isRoot, nil
}
