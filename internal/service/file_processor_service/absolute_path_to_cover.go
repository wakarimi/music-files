package file_processor_service

import (
	"github.com/jmoiron/sqlx"
	"path/filepath"
)

func (s *Service) AbsolutePathToCover(tx *sqlx.Tx, coverId int) (absolutePath string, err error) {
	cover, err := s.CoverService.GetCover(tx, coverId)
	if err != nil {
		return "", err
	}

	absolutePathToDir, err := s.DirService.AbsolutePath(tx, cover.DirId)
	if err != nil {
		return "", err
	}

	absolutePath = filepath.Join(absolutePathToDir, cover.Filename)

	return absolutePath, nil
}
