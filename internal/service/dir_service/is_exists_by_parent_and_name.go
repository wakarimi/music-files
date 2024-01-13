package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s Service) IsExistsByParentAndName(tx *sqlx.Tx, parentID *int, name string) (exists bool, err error) {
	log.Debug().Interface("parentID", parentID).Str("name", name).Msg("Checking directory existence")

	exists, err = s.dirRepo.IsExistsByParentAndName(tx, parentID, name)
	if err != nil {
		log.Error().Err(err).Interface("parentID", parentID).Str("name", name).Msg("Failed to check directory existence")
		return false, err
	}

	log.Debug().Interface("parentID", parentID).Str("name", name).Bool("exists", exists).Msg("Directory existence checked")
	return exists, nil
}
