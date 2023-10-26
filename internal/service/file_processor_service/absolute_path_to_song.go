package file_processor_service

import (
	"github.com/jmoiron/sqlx"
	"path/filepath"
)

func (s *Service) AbsolutePathToSong(tx *sqlx.Tx, songId int) (absolutePath string, err error) {
	song, err := s.SongService.GetSong(tx, songId)
	if err != nil {
		return "", err
	}

	absolutePathToDir, err := s.DirService.AbsolutePath(tx, song.DirId)
	if err != nil {
		return "", err
	}

	absolutePath = filepath.Join(absolutePathToDir, song.Filename)

	return absolutePath, nil
}
