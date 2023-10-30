package file_processor_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"path/filepath"
)

func (s *Service) AbsolutePathToCover(tx *sqlx.Tx, coverId int) (absolutePath string, err error) {
	log.Debug().Int("coverId", coverId).Msg("Calculating absolute path to cover")

	cover, err := s.CoverService.GetCover(tx, coverId)
	if err != nil {
		log.Debug().Int("coverId", coverId).Msg("Failed to get cover")
		return "", err
	}

	absolutePathToDir, err := s.DirService.AbsolutePath(tx, cover.DirId)
	if err != nil {
		log.Debug().Int("dirId", cover.DirId).Msg("Failed to calculate absolute path to directory")
		return "", err
	}

	absolutePath = filepath.Join(absolutePathToDir, cover.Filename)

	log.Debug().Int("coverId", coverId).Str("absolutePath", absolutePath).Msg("Calculating absolute path to cover")
	return absolutePath, nil
}
