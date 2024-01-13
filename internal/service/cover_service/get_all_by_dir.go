package cover_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/cover"
)

func (s Service) GetAllByDir(tx *sqlx.Tx, dirID int) (covers []cover.Cover, err error) {
	log.Debug().Int("dirId", dirID).Msg("Fetching cover")

	covers, err = s.coverRepo.ReadAllByDir(tx, dirID)
	if err != nil {
		log.Error().Int("dirId", dirID).Err(err).Msg("Failed to fetch all cover")
		return make([]cover.Cover, 0), err
	}

	log.Debug().Int("dirId", dirID).Int("countOfcover", len(covers)).Msg("All cover files fetched successfully")
	return covers, nil
}
