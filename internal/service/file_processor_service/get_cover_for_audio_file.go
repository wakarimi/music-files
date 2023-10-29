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
		return models.Cover{}, err
	}

	dir, err := s.DirService.GetDir(tx, audioFile.DirId)
	if err != nil {
		return models.Cover{}, err
	}

	covers, err := s.CoverService.GetAllByDir(tx, dir.DirId)
	if err != nil {
		return models.Cover{}, err
	}

	for (len(covers) == 0) && (dir.ParentDirId != nil) {
		dir, err = s.DirService.GetDir(tx, audioFile.DirId)
		if err != nil {
			return models.Cover{}, err
		}

		covers, err = s.CoverService.GetAllByDir(tx, dir.DirId)
		if err != nil {
			return models.Cover{}, err
		}
	}

	if len(covers) == 0 {
		err = errors.NotFound{Resource: fmt.Sprintf("cover for audioFile with id=%d", audioFileId)}
		return models.Cover{}, err
	}

	log.Debug().Int("coverId", covers[0].CoverId).Msg("Cover for audioFile got successfully")
	return covers[0], nil
}
