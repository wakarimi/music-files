package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/directory"
	"strings"
)

func (s Service) ContainedRoots(tx *sqlx.Tx, path string) ([]directory.Directory, error) {
	log.Debug().Str("path", path).Msg("Finding root child folders for a specified path")

	rootDirs, err := s.dirRepo.ReadRoots(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read roots")
		return make([]directory.Directory, 0), err
	}

	childRoots := make([]directory.Directory, 0)
	for _, rootDir := range rootDirs {
		if strings.HasPrefix(rootDir.Name, path) {
			childRoots = append(childRoots, rootDir)
		}
	}

	log.Debug().Int("count", len(childRoots)).Msg("Child roots read")
	return childRoots, nil
}
