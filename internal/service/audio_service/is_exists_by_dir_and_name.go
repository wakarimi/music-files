package audio_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s Service) IsExistsByDirAndName(tx *sqlx.Tx, dirID int, name string) (exists bool, err error) {
	log.Debug().Int("dirId", dirID).Str("name", name).Msg("Checking audio file existence")

	exists, err = s.audioRepo.IsExistsByDirAndName(tx, dirID, name)
	if err != nil {
		log.Debug().Int("dirId", dirID).Str("name", name).Msg("Failed to check audio file existence")
		return false, err
	}

	log.Debug().Int("dirId", dirID).Str("name", name).Bool("exists", exists).Msg("Audio file existence checked successfully")
	return exists, nil
}
