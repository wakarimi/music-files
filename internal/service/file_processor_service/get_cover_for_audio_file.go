package file_processor_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/models"
)

func (s *Service) GetCoverForAudioFile(tx *sqlx.Tx, audioFileId int) (cover models.Cover, err error) {
	log.Debug().Int("audioFileId", audioFileId).Msg("Getting cover for audioFile")

	audioFile, err := s.AudioFileService.GetAudioFile(tx, audioFileId)
	if err != nil {
		log.Error().Err(err).Int("audioFileId", audioFileId).Msg("Failed to get audio file")
		return models.Cover{}, err
	}

	dir, err := s.DirService.GetDir(tx, audioFile.DirId)
	if err != nil {
		log.Error().Err(err).Int("dirId", audioFile.DirId).Msg("Failed to get audio file's directory")
		return models.Cover{}, err
	}

	covers, err := s.CoverService.GetAllByDir(tx, dir.DirId)
	if err != nil {
		log.Error().Err(err).Int("dirId", dir.DirId).Msg("Failed to get covers in directory")
		return models.Cover{}, err
	}

	for (len(covers) == 0) && (dir.ParentDirId != nil) {
		dir, err = s.DirService.GetDir(tx, *dir.ParentDirId)
		if err != nil {
			log.Error().Err(err).Int("dirId", audioFile.DirId).Msg("Failed to get subdirectories")
			return models.Cover{}, err
		}

		covers, err = s.CoverService.GetAllByDir(tx, dir.DirId)
		if err != nil {
			log.Error().Err(err).Int("dirId", dir.DirId).Msg("Failed to get covers in directory")
			return models.Cover{}, err
		}
	}

	if len(covers) == 0 {
		log.Error().Err(err).Int("audioFileId", audioFileId).Msg("Cover foa audio file not found")
		return models.Cover{}, errors.NotFound{Resource: fmt.Sprintf("cover for audio_file with id=%d", audioFileId)}
	}

	log.Debug().Int("coverId", covers[0].CoverId).Msg("Cover for audioFile got successfully")
	return covers[0], nil
}
