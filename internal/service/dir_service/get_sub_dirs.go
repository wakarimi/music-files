package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/directory"
)

func (s Service) GetSubDirs(tx *sqlx.Tx, dirID int) (subDirs []directory.Directory, err error) {
	log.Debug().Int("dirId", dirID).Msg("Getting subdirectories")

	subDirs, err = s.dirRepo.ReadSubDirs(tx, dirID)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirID).Msg("Failed to get subdirectories")
		return make([]directory.Directory, 0), err
	}

	log.Debug().Int("count", len(subDirs)).Msg("Subdirectories got")
	return subDirs, nil
}
