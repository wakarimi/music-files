package cover_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s Service) DeleteAllByDir(tx *sqlx.Tx, dirID int) (err error) {
	log.Debug().Int("dirId", dirID).Msg("Deleting audios by directory")

	err = s.coverRepo.DeleteByDir(tx, dirID)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirID).Msg("Failed to delete audios by directory")
		return err
	}

	log.Debug().Int("dirId", dirID).Msg("Audios by directory deleted")
	return nil
}
