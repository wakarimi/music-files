package file_processor_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/models"
)

func (s *Service) GetCoverForSong(tx *sqlx.Tx, songId int) (cover models.Cover, err error) {
	log.Debug().Int("songId", songId).Msg("Getting cover for song")

	song, err := s.SongService.GetSong(tx, songId)
	if err != nil {
		return models.Cover{}, err
	}

	dir, err := s.DirService.GetDir(tx, song.DirId)
	if err != nil {
		return models.Cover{}, err
	}

	covers, err := s.CoverService.GetAllByDir(tx, dir.DirId)
	if err != nil {
		return models.Cover{}, err
	}

	for (len(covers) == 0) && (dir.ParentDirId != nil) {
		dir, err = s.DirService.GetDir(tx, song.DirId)
		if err != nil {
			return models.Cover{}, err
		}

		covers, err = s.CoverService.GetAllByDir(tx, dir.DirId)
		if err != nil {
			return models.Cover{}, err
		}
	}

	if len(covers) == 0 {
		err = errors.NotFound{Resource: fmt.Sprintf("cover for song with id=%d", songId)}
		return models.Cover{}, err
	}

	log.Debug().Int("coverId", covers[0].CoverId).Msg("Cover for song got successfully")
	return covers[0], nil
}
