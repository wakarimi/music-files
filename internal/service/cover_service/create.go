package cover_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/cover"
)

func (s Service) Create(tx *sqlx.Tx, coverToCreate cover.Cover) (int, error) {
	log.Debug().Interface("coverToCreate", coverToCreate).Msg("Creating new cover file")

	coverFileID, err := s.coverRepo.Create(tx, coverToCreate)
	if err != nil {
		log.Error().Err(err).Interface("coverToCreate", coverToCreate).Msg("Failed to create new cover file")
		return 0, err
	}

	log.Debug().Interface("coverToCreate", coverToCreate).Msg("cover file created successfully")
	return coverFileID, nil
}
