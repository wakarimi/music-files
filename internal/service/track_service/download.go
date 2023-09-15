package track_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

func (s *Service) Download(tx *sqlx.Tx, trackId int) (absolutePath string, err error) {
	track, err := s.TrackRepo.ReadTx(tx, trackId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read track")
		return "", err
	}
	log.Debug().Str("relativePath", track.RelativePath).Msg("Track read successfully")

	dir, err := s.DirRepo.Read(track.DirId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read dir")
		return "", err
	}
	log.Debug().Str("path", dir.Path).Msg("Dir read successfully")

	absolutePath = filepath.Join(dir.Path, track.RelativePath, track.Filename)
	file, err := os.Open(absolutePath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open track file")
		return "", err
	}
	defer file.Close()
	log.Debug().Str("filename", file.Name()).Msg("File read successfully")

	return absolutePath, nil
}
