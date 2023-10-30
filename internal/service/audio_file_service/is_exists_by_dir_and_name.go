package audio_file_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s *Service) IsExistsByDirAndName(tx *sqlx.Tx, dirId int, name string) (exists bool, err error) {
	log.Debug().Int("dirId", dirId).Str("name", name).Msg("Checking audio file existence")

	exists, err = s.AudioFileRepo.IsExistsByDirAndName(tx, dirId, name)
	if err != nil {
		log.Debug().Int("dirId", dirId).Str("name", name).Msg("Failed to check audio file existence")
		return false, err
	}

	log.Debug().Int("dirId", dirId).Str("name", name).Bool("exists", exists).Msg("Audio file existence checked successfully")
	return exists, nil
}
