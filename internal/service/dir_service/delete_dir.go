package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s *Service) DeleteDir(tx *sqlx.Tx, dirId int) (err error) {
	log.Debug().Int("dirId", dirId).Msg("Deleting directory with files from database")

	subDirs, err := s.DirRepo.ReadSubDirs(tx, dirId)
	if err != nil {
		log.Warn().Err(err).Int("dirId", dirId).Msg("Failed to get subdirectories")
		return err
	}

	for _, subDir := range subDirs {
		err := s.DeleteDir(tx, subDir.DirId)
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

	songs, err := s.SongService.GetAllByDir(tx, dirId)
	if err != nil {
		return err
	}

	for _, song := range songs {
		err = s.SongService.Delete(tx, song.SongId)
		if err != nil {
			return err
		}
	}

	covers, err := s.CoverService.GetAllByDir(tx, dirId)
	if err != nil {
		return err
	}

	for _, cover := range covers {
		err = s.CoverService.Delete(tx, cover.CoverId)
		if err != nil {
			return err
		}
	}

	log.Debug().Int("dirId", dirId).Msg("Files deleted from directory successfully")
	return nil
}
