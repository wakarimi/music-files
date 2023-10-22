package dir_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
)

func (s *Service) UntrackRoot(tx *sqlx.Tx, dirId int) (err error) {
	log.Debug().Int("dirId", dirId).Msg("Removing directories from tracked")

	exists, err := s.DirRepo.IsExists(tx, dirId)
	if err != nil {
		log.Warn().Err(err).Int("dirId", dirId).Msg("Failed to check directory existing")
		return err
	}
	if !exists {
		err = errors.NotFound{Resource: fmt.Sprintf("directory with id=%d", dirId)}
		log.Warn().Err(err).Int("dirId", dirId).Msg("Directory not found in database")
		return err
	}

	dir, err := s.DirRepo.Read(tx, dirId)
	if err != nil {
		log.Warn().Err(err).Int("dirId", dirId).Msg("Failed to get directory")
		return err
	}
	if dir.ParentDirId != nil {
		err = errors.BadRequest{Message: fmt.Sprintf("directory with id=%d is not root", dirId)}
		log.Warn().Err(err).Int("dirId", dirId).Msg("Directory is not root")
		return err
	}

	subDirs, err := s.DirRepo.ReadSubDirs(tx, dirId)
	if err != nil {
		log.Warn().Err(err).Int("dirId", dirId).Msg("Failed to get subdirectories")
		return err
	}
	for _, subDir := range subDirs {
		err := s.deleteDir(tx, subDir.DirId)
		if err != nil {
			log.Warn().Err(err).Int("subDirId", subDir.DirId).Msg("Failed to delete subdirectory from database")
			return err
		}
	}

	err = s.deleteDir(tx, dirId)
	if err != nil {
		log.Warn().Err(err).Int("dirId", dirId).Msg("Failed to delete directory from database")
	}

	log.Debug().Int("dirId", dirId).Msg("Directory removed from tracked")
	return nil
}

func (s *Service) deleteDir(tx *sqlx.Tx, dirId int) (err error) {
	log.Debug().Int("dirId", dirId).Msg("Deleting directory with files from database")

	subDirs, err := s.DirRepo.ReadSubDirs(tx, dirId)
	if err != nil {
		log.Warn().Err(err).Int("dirId", dirId).Msg("Failed to get subdirectories")
		return err
	}

	for _, subDir := range subDirs {
		err := s.deleteDir(tx, subDir.DirId)
		if err != nil {
			log.Warn().Err(err).Int("subDirId", subDir.DirId).Msg("Failed to delete subdirectory from database")
			return err
		}
	}

	err = s.deleteContentFiles(tx, dirId)
	if err != nil {
		log.Warn().Err(err).Int("dirId", dirId).Msg("Failed to delete files from database")
		return err
	}

	err = s.DirRepo.Delete(tx, dirId)
	if err != nil {
		log.Warn().Err(err).Int("dirId", dirId).Msg("Failed to delete directory from database")
	}

	log.Debug().Int("dirId", dirId).Msg("Directory with files deleted from database successfully")
	return nil
}

func (s *Service) deleteContentFiles(tx *sqlx.Tx, dirId int) (err error) {
	log.Debug().Int("dirId", dirId).Msg("Deleting files in directory")

	// TODO: Delete songs

	// TODO: Delete covers

	log.Debug().Int("dirId", dirId).Msg("Files deleted from directory successfully")
	return nil
}
