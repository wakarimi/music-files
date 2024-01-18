package cover_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/cover"
)

func (s Service) Get(tx *sqlx.Tx, coverID int) (coverFile cover.Cover, err error) {
	log.Debug().Int("coverId", coverID).Msg("Getting cover file")

	coverFile, err = s.coverRepo.Read(tx, coverID)
	if err != nil {
		log.Error().Err(err).Int("coverId", coverID).Msg("Failed to fetch cover file")
		return cover.Cover{}, err
	}

	log.Debug().Int("coverId", coverID).Msg("Cover file got successfully")
	return coverFile, nil
}
