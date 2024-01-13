package cover_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/cover"
)

func (s Service) Update(tx *sqlx.Tx, coverID int, coverToUpdate cover.Cover) (err error) {
	log.Debug().Int("coverId", coverID).Msg("Updating cover file")

	err = s.coverRepo.Update(tx, coverID, coverToUpdate)
	if err != nil {
		log.Debug().Int("coverId", coverID).Msg("Failed to update cover file")
		return err
	}

	log.Debug().Int("coverId", coverID).Msg("Cover file updated successfully")
	return nil
}
