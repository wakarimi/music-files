package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s *Service) DeleteDir(tx *sqlx.Tx, dirId int) (err error) {
	log.Debug().Int("dirId", dirId).Msg("Deleting directory from database")

	subDirs, err := s.DirRepo.ReadSubDirs(tx, dirId)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Msg("Failed to get subdirectories")
		return err
	}
	log.Debug().Int("dirId", dirId).Int("countOfSubDirs", len(subDirs)).Msg("Subdirectories taken")

	for _, subDir := range subDirs {
		err := s.DeleteDir(tx, subDir.DirId)
		if err != nil {
			log.Error().Err(err).Int("subDirId", subDir.DirId).Msg("Failed to delete subdirectory from database")
			return err
		}
	}

	err = s.deleteContent(tx, dirId)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Msg("Failed to delete files from database")
		return err
	}

	err = s.DirRepo.Delete(tx, dirId)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Msg("Failed to delete directory from database")
		return err
	}

	log.Debug().Int("dirId", dirId).Msg("Directory with files deleted from database successfully")
	return nil
}

func (s *Service) deleteContent(tx *sqlx.Tx, dirId int) (err error) {
	log.Debug().Int("dirId", dirId).Msg("Deleting files in directory")

	audioFiles, err := s.AudioFileService.GetAllByDir(tx, dirId)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Msg("Failed to get directory's audio files")
		return err
	}
	log.Debug().Int("dirId", dirId).Int("countOfAudioFileInDir", len(audioFiles)).Msg("Received audioFiles that need to be deleted")

	for _, audioFile := range audioFiles {
		err = s.AudioFileService.Delete(tx, audioFile.AudioFileId)
		if err != nil {
			log.Error().Err(err).Int("audioFileId", audioFile.AudioFileId).Msg("Failed to delete audioFile from database")
			return err
		}
	}

	covers, err := s.CoverService.GetAllByDir(tx, dirId)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Msg("Failed to get directory's covers")
		return err
	}
	log.Debug().Int("dirId", dirId).Int("countOfCoverInDir", len(audioFiles)).Msg("Received covers that need to be deleted")

	for _, cover := range covers {
		err = s.CoverService.Delete(tx, cover.CoverId)
		if err != nil {
			log.Error().Err(err).Int("coverId", cover.CoverId).Msg("Failed to delete cover from database")
			return err
		}
	}

	log.Debug().Int("dirId", dirId).Msg("Files deleted from directory successfully")
	return nil
}
