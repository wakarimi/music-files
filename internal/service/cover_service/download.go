package cover_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

func (s *Service) Download(tx *sqlx.Tx, coverId int) (absolutePath string, err error) {
	cover, err := s.CoverRepo.ReadTx(tx, coverId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read cover")
		return "", err
	}
	log.Debug().Str("relativePath", cover.RelativePath).Msg("Cover read successfully")

	dir, err := s.DirRepo.Read(cover.DirId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read dir")
		return "", err
	}
	log.Debug().Str("path", dir.Path).Msg("Dir read successfully")

	absolutePath = filepath.Join(dir.Path, cover.RelativePath, cover.Filename)
	file, err := os.Open(absolutePath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open cover file")
		return "", err
	}
	defer file.Close()
	log.Debug().Str("filename", file.Name()).Msg("File read successfully")

	return absolutePath, nil
}
