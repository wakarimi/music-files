package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"path/filepath"
)

func (s Service) CalcAbsolutePath(tx *sqlx.Tx, dirID int) (absolutePath string, err error) {
	log.Debug().Int("dirId", dirID).Msg("Calculating absolute path")

	var parts []string
	currentDir, err := s.dirRepo.Read(tx, dirID)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirID).Msg("Failed to fetch directory from database")
		return "", err
	}

	parts = append([]string{currentDir.Name}, parts...)
	for currentDir.ParentDirID != nil {
		currentDir, err = s.dirRepo.Read(tx, *currentDir.ParentDirID)
		if err != nil {
			log.Error().Err(err).Int("dirId", dirID).Int("currentDirId", currentDir.ID).Msg("Failed to fetch parent directory from database")
			return "", err
		}
		parts = append([]string{currentDir.Name}, parts...)
	}
	absolutePath = filepath.Join(parts...)

	log.Debug().Int("dirId", dirID).Str("absolutePath", absolutePath).Msg("Absolute path calculated successfully")
	return absolutePath, nil
}
