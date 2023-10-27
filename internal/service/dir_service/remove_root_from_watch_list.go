package dir_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
)

func (s *Service) RemoveRootFromWatchList(tx *sqlx.Tx, dirId int) (err error) {
	log.Debug().Int("dirId", dirId).Msg("Removing root directory from watch list")

	exists, err := s.DirRepo.IsExists(tx, dirId)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Msg("Failed to check directory existing")
		return err
	}
	if !exists {
		err = errors.NotFound{Resource: fmt.Sprintf("directory with id=%d", dirId)}
		log.Error().Err(err).Int("dirId", dirId).Msg("Directory not found in database")
		return err
	}

	dir, err := s.DirRepo.Read(tx, dirId)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Msg("Failed to get directory")
		return err
	}
	if dir.ParentDirId != nil {
		err = errors.BadRequest{Message: fmt.Sprintf("directory with id=%d is not root", dirId)}
		log.Error().Err(err).Int("dirId", dirId).Msg("Directory is not root")
		return err
	}

	err = s.DeleteDir(tx, dirId)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Msg("Failed to delete directory from database")
	}

	log.Debug().Int("dirId", dirId).Msg("Directory removed from watch list")
	return nil
}
