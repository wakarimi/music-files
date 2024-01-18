package cover_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s Service) IsExists(tx *sqlx.Tx, coverID int) (exists bool, err error) {
	log.Debug().Int("coverId", coverID).Msg("Checking cover file existence")

	exists, err = s.coverRepo.IsExists(tx, coverID)
	if err != nil {
		log.Debug().Int("coverId", coverID).Msg("Failed to check cover file existence")
		return false, err
	}

	log.Debug().Int("coverId", coverID).Bool("exists", exists).Msg("Cover file existence checked successfully")
	return exists, nil
}
