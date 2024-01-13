package cover_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/cover"
)

func (s Service) GetByDirAndName(tx *sqlx.Tx, dirID int, name string) (coverFile cover.Cover, err error) {
	log.Debug().Int("dirId", dirID).Str("name", name).Msg("Getting cover file")

	coverFile, err = s.coverRepo.ReadByDirAndName(tx, dirID, name)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirID).Str("name", name).Msg("Failed to fetch cover file")
		return cover.Cover{}, err
	}

	log.Debug().Int("dirId", dirID).Str("name", name).Msg("cover file got successfully")
	return coverFile, nil
}
