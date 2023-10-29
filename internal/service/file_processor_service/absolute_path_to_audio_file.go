package file_processor_service

import (
	"github.com/jmoiron/sqlx"
	"path/filepath"
)

func (s *Service) AbsolutePathToAudioFile(tx *sqlx.Tx, audioFileId int) (absolutePath string, err error) {
	audioFile, err := s.AudioFileService.GetAudioFile(tx, audioFileId)
	if err != nil {
		return "", err
	}

	absolutePathToDir, err := s.DirService.AbsolutePath(tx, audioFile.DirId)
	if err != nil {
		return "", err
	}

	absolutePath = filepath.Join(absolutePathToDir, audioFile.Filename)

	return absolutePath, nil
}
