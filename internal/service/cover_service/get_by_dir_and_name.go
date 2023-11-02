package cover_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/model"
)

func (s *Service) GetByDirAndName(tx *sqlx.Tx, dirId int, name string) (cover model.Cover, err error) {
	log.Debug().Int("dirId", dirId).Str("name", name).Msg("Getting cover")

	exists, err := s.CoverRepo.IsExistsByDirAndName(tx, dirId, name)
	if err != nil {
		log.Error().Int("dirId", dirId).Str("name", name).Msg("Failed to check cover existence")
		return model.Cover{}, err
	}
	if !exists {
		log.Error().Int("dirId", dirId).Str("name", name).Msg("Cover not found")
		return model.Cover{}, errors.NotFound{Resource: fmt.Sprintf("cover with dirId=%d and name=%s in database", dirId, name)}
	}

	cover, err = s.CoverRepo.ReadByDirAndName(tx, dirId, name)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Str("name", name).Msg("Failed to fetch cover")
		return model.Cover{}, err
	}

	log.Debug().Int("dirId", dirId).Str("name", name).Msg("Cover got successfully")
	return cover, nil
}
