package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"path/filepath"
)

func (s *Service) AbsolutePath(tx *sqlx.Tx, dirId int) (absolutePath string, err error) {
	log.Debug().Int("dirId", dirId).Msg("Calculating absolute path")

	var parts []string
	currentDir, err := s.DirRepo.Read(tx, dirId)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Msg("Failed to fetch directory from database")
		return "", err
	}

	parts = append([]string{currentDir.Name}, parts...)
	for currentDir.ParentDirId != nil {
		currentDir, err = s.DirRepo.Read(tx, *currentDir.ParentDirId)
		if err != nil {
			log.Error().Err(err).Int("dirId", dirId).Int("currentDirId", currentDir.DirId).Msg("Failed to fetch parent directory from database")
			return "", err
		}
		parts = append([]string{currentDir.Name}, parts...)
	}
	absolutePath = filepath.Join(parts...)

	log.Debug().Int("dirId", dirId).Str("absolutePath", absolutePath).Msg("Absolute path calculated successfully")
	return absolutePath, nil
}
