package dir_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/models"
	"music-files/internal/utils"
	"path/filepath"
	"strings"
)

func (s *Service) AddRootToWatchList(tx *sqlx.Tx, dir models.Directory) (createdDir models.Directory, err error) {
	log.Debug().Str("path", dir.Name).Msg("Adding a new directory tracking")

	dir.Name = filepath.Clean(dir.Name)

	directoryExists, err := utils.IsDirectoryExistsOnDisk(dir.Name)
	if err != nil {
		log.Error().Err(err).Str("path", dir.Name).Msg("Failed to check directory on disk")
		return models.Directory{}, err
	}
	if !directoryExists {
		log.Error().Err(err).Str("path", dir.Name).Msg("Directory nod found on disk")
		return models.Directory{}, errors.NotFound{Resource: "directory on disk"}
	}

	alreadyTracked, err := s.isAlreadyInWatchList(tx, dir.Name)
	if err != nil {
		log.Error().Err(err).Str("path", dir.Name).Msg("Couldn't check if the directory is already in the watch list")
		return models.Directory{}, err
	}
	if alreadyTracked {
		err = fmt.Errorf("directory with path=%s already in watch list", dir.Name)
		log.Error().Err(err).Str("path", dir.Name).Msg("Directory already in watch list")
		return models.Directory{}, errors.Conflict{Message: err.Error()}
	}

	createdDirId, err := s.DirRepo.Create(tx, dir)
	if err != nil {
		log.Error().Err(err).Str("path", dir.Name).Msg("Failed to add directory to watch list")
		return models.Directory{}, err
	}
	createdDir, err = s.DirRepo.Read(tx, createdDirId)
	if err != nil {
		log.Error().Err(err).Int("dirId", createdDir.DirId).Msg("Failed to read added directory")
		return models.Directory{}, err
	}

	containedRoots, err := s.getContainedRoots(tx, dir.Name)
	if err != nil {
		log.Error().Err(err).Str("path", dir.Name).Msg("Failed to get root subdirectories")
		return models.Directory{}, err
	}
	err = s.addIntermediateDirs(tx, containedRoots)
	if err != nil {
		log.Error().Err(err).Str("path", dir.Name).Msg("Failed to add intermediate dirs")
		return models.Directory{}, err
	}

	log.Debug().Str("path", dir.Name).Msg("Directory added to tracked")
	return createdDir, nil
}

func (s *Service) isAlreadyInWatchList(tx *sqlx.Tx, absolutePath string) (alreadyTracked bool, err error) {
	trackedDirs, err := s.DirRepo.ReadAll(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all dirs")
		return false, err
	}

	for _, trackedDir := range trackedDirs {
		trackedDirPath, err := s.AbsolutePath(tx, trackedDir.DirId)
		if err != nil {
			log.Error().Err(err).Int("trackedDirId", trackedDir.DirId).Msg("Failed to calculate absolute path")
			return false, err
		}
		if strings.Contains(absolutePath, trackedDirPath) {
			return true, nil
		}
	}

	return false, nil
}

func (s *Service) getContainedRoots(tx *sqlx.Tx, absolutePath string) (containedRoots []models.Directory, err error) {
	log.Debug().Str("path", absolutePath).Msg("Getting root directories that are subdirectories of the current one")

	containedRoots = make([]models.Directory, 0)

	roots, err := s.DirRepo.ReadRoots(tx)
	if err != nil {
		log.Error().Err(err).Str("absolutePath", absolutePath).Msg("Failed to get root dirs")
		return make([]models.Directory, 0), err
	}

	for _, root := range roots {
		if strings.Contains(root.Name, absolutePath) && root.Name != absolutePath {
			containedRoots = append(containedRoots, root)
		}
	}

	log.Debug().Str("currentDirPath", absolutePath).Msg("Root directories that are subdirectories of the current one are got")
	return containedRoots, nil
}

func (s *Service) addIntermediateDirs(tx *sqlx.Tx, containedRoots []models.Directory) (err error) {
	log.Debug().Int("countOfRootSubDirs", len(containedRoots)).Msg("Adding paths between the current directory and subdirectory root")

	for _, containedRoot := range containedRoots {
		log.Info().Str("path", containedRoot.Name).Msg("The directory is nested inside another and will be absorbed")

		prevDirId := containedRoot.DirId
		currentPath, folder := filepath.Split(containedRoot.Name)
		currentPath = filepath.Clean(currentPath)
		alreadyTracked, err := s.isEqualDirTracked(tx, currentPath)
		if err != nil {
			log.Error().Err(err).Str("currentPath", currentPath).Msg("Failed to check equals dir")
			return err
		}
		for !alreadyTracked {
			newDirId, err := s.DirRepo.Create(tx, models.Directory{
				ParentDirId: nil,
				Name:        currentPath,
			})
			if err != nil {
				log.Error().Err(err).Str("currentPath", currentPath).Msg("Failed to create new dir in database")
				return err
			}
			err = s.DirRepo.Update(tx, prevDirId, models.Directory{
				ParentDirId: &newDirId,
				Name:        folder,
			})
			if err != nil {
				log.Error().Err(err).Int("subDirId", prevDirId).Msg("Failed to update subdirectory")
				return err
			}
			prevDirId = newDirId
			currentPath, folder = filepath.Split(currentPath)
			currentPath = filepath.Clean(currentPath)
			alreadyTracked, err = s.isEqualDirTracked(tx, currentPath)
		}
		topDir, err := s.getEqualTrackedDir(tx, currentPath)
		if err != nil {
			log.Error().Err(err).Str("currentPath", currentPath).Msg("Failed to get equal dir from watch list")
			return err
		}
		err = s.DirRepo.Update(tx, prevDirId, models.Directory{
			ParentDirId: &topDir.DirId,
			Name:        folder,
		})
		if err != nil {
			log.Error().Err(err).Int("subDirId", prevDirId).Msg("Failed to update subdirectory")
			return err
		}
	}

	log.Debug().Msg("Intermediate paths added successfully")
	return nil
}

func (s *Service) isEqualDirTracked(tx *sqlx.Tx, absolutePath string) (equalDirTracked bool, err error) {
	log.Debug().Str("path", absolutePath).Msg("Checking to see if a directory is already being in watch list")

	trackedDirs, err := s.DirRepo.ReadAll(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch all dirs")
		return false, err
	}

	for _, trackedDir := range trackedDirs {
		trackedDirPath, err := s.AbsolutePath(tx, trackedDir.DirId)
		if err != nil {
			log.Error().Err(err).Int("trackedDirId", trackedDir.DirId).Msg("Failed to calculate absolute path")
			return false, err
		}
		if absolutePath == trackedDirPath {
			log.Debug().Str("path", absolutePath).Msg("The directory is already on the watch list")
			return true, nil
		}
	}

	log.Debug().Str("path", absolutePath).Msg("The directory is not yet on the watch list")
	return false, nil
}

func (s *Service) getEqualTrackedDir(tx *sqlx.Tx, absolutePath string) (equalTrackedDir models.Directory, err error) {
	log.Debug().Msg("Getting a similar directory from the database")

	trackedDirs, err := s.DirRepo.ReadAll(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch all dirs")
		return models.Directory{}, err
	}

	for _, trackedDir := range trackedDirs {
		trackedDirPath, err := s.AbsolutePath(tx, trackedDir.DirId)
		if err != nil {
			log.Error().Err(err).Int("trackedDirId", trackedDir.DirId).Msg("Failed to calculate absolute path")
			return models.Directory{}, err
		}
		if absolutePath == trackedDirPath {
			log.Debug().Str("path", absolutePath).Msg("Similar directory found")
			return trackedDir, nil
		}
	}

	err = fmt.Errorf("similar directory not found, path=%s", absolutePath)
	log.Error().Err(err).Msg("Similar directory not found")
	return models.Directory{}, err
}
