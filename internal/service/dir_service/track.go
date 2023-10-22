package dir_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/models"
	"os"
	"path/filepath"
	"strings"
)

func (s *Service) Track(tx *sqlx.Tx, dir models.Directory) (createdDir models.Directory, err error) {
	log.Debug().Str("path", dir.Name).Msg("Adding a new directory tracking")

	dir.Name = filepath.Clean(dir.Name)

	directoryExists, err := s.isDirectoryExistsOnDisk(dir)
	if err != nil {
		log.Warn().Err(err).Str("path", dir.Name).Msg("Failed to check directory on disk")
		return models.Directory{}, err
	}
	if !directoryExists {
		log.Info().Err(err).Str("path", dir.Name).Msg("Directory nod found on disk")
		return models.Directory{}, errors.NotFound{Resource: "directory on disk"}
	}

	alreadyTracked, err := s.isAlreadyTracked(tx, dir.Name)
	if err != nil {
		log.Warn().Err(err).Str("path", dir.Name).Msg("Failed to check whether directory was added")
		return models.Directory{}, err
	}
	if alreadyTracked {
		err = fmt.Errorf("directory with path=%s already tracked", dir.Name)
		log.Info().Err(err).Str("path", dir.Name).Msg("Directory already tracked")
		return models.Directory{}, errors.Conflict{Message: err.Error()}
	}

	createdDirId, err := s.DirRepo.Create(tx, dir)
	if err != nil {
		log.Warn().Err(err).Str("path", dir.Name).Msg("Failed to track directory")
		return models.Directory{}, err
	}
	createdDir, err = s.DirRepo.Read(tx, createdDirId)
	if err != nil {
		log.Warn().Err(err).Int("dirId", createdDir.DirId).Msg("Failed to read added directory")
		return models.Directory{}, err
	}

	containedRoots, err := s.getContainedRoots(tx, dir.Name)
	if err != nil {
		log.Warn().Err(err).Str("path", dir.Name).Msg("Failed to get contained roots")
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

func (s *Service) isDirectoryExistsOnDisk(dir models.Directory) (directoryExists bool, err error) {
	_, err = os.Stat(dir.Name)
	if err != nil {
		log.Warn().Err(err).Str("path", dir.Name).Msg("Failed to check directory on disk")
		return false, err
	}
	if os.IsNotExist(err) {
		log.Info().Err(err).Str("path", dir.Name).Msg("Directory on disk not found")
		return false, errors.NotFound{Resource: "directory on disk"}
	} else if err != nil {
		log.Warn().Err(err).Str("path", dir.Name).Msg("Unknown error when checking for existence")
		return false, err
	}

	fileInfo, err := os.Stat(dir.Name)
	if err != nil {
		log.Warn().Err(err).Str("path", dir.Name).Msg("Failed to get object info")
		return false, err
	}
	if !fileInfo.IsDir() {
		log.Info().Err(err).Str("filepath", dir.Name).Msg("Trying to add a file instead of a folder")
		return false, err
	}

	return true, nil
}

func (s *Service) isAlreadyTracked(tx *sqlx.Tx, absolutePath string) (alreadyTracked bool, err error) {
	trackedDirs, err := s.DirRepo.ReadAll(tx)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to fetch all dirs")
		return false, err
	}

	for _, trackedDir := range trackedDirs {
		trackedDirPath, err := s.AbsolutePath(tx, trackedDir.DirId)
		if err != nil {
			log.Warn().Err(err).Int("trackedDirId", trackedDir.DirId).Msg("Failed to calculate absolute path")
			return false, err
		}
		if strings.Contains(absolutePath, trackedDirPath) {
			return true, nil
		}
	}

	return false, nil
}

func (s *Service) getContainedRoots(tx *sqlx.Tx, absolutePath string) (containedRoots []models.Directory, err error) {
	containedRoots = make([]models.Directory, 0)

	roots, err := s.DirRepo.ReadRoots(tx)
	if err != nil {
		log.Warn().Err(err).Str("absolutePath", absolutePath).Msg("Failed to get root dirs")
		return make([]models.Directory, 0), err
	}

	for _, root := range roots {
		if strings.Contains(root.Name, absolutePath) && root.Name != absolutePath {
			containedRoots = append(containedRoots, root)
		}
	}
	return containedRoots, nil
}

func (s *Service) addIntermediateDirs(tx *sqlx.Tx, containedRoots []models.Directory) (err error) {
	for _, containedRoot := range containedRoots {
		log.Info().Str("path", containedRoot.Name).Msg("The directory is nested inside another and will be absorbed")

		prevDirId := containedRoot.DirId
		currentPath, folder := filepath.Split(containedRoot.Name)
		currentPath = filepath.Clean(currentPath)
		alreadyTracked, err := s.isEqualDirTracked(tx, currentPath)
		if err != nil {
			log.Warn().Err(err).Str("currentPath", currentPath).Msg("Failed to check equals dir")
			return err
		}
		for !alreadyTracked {
			newDirId, err := s.DirRepo.Create(tx, models.Directory{
				ParentDirId: nil,
				Name:        currentPath,
			})
			if err != nil {
				log.Warn().Err(err).Str("currentPath", currentPath).Msg("Failed to create new dir in database")
				return err
			}
			err = s.DirRepo.Update(tx, prevDirId, models.Directory{
				ParentDirId: &newDirId,
				Name:        folder,
			})
			if err != nil {
				log.Warn().Err(err).Int("subDirId", prevDirId).Msg("Failed to update subdirectory")
				return err
			}
			prevDirId = newDirId
			currentPath, folder = filepath.Split(currentPath)
			currentPath = filepath.Clean(currentPath)
			alreadyTracked, err = s.isEqualDirTracked(tx, currentPath)
		}
		topDir, err := s.getEqualTrackedDir(tx, currentPath)
		if err != nil {
			log.Warn().Err(err).Str("currentPath", currentPath).Msg("Failed to get equalTrackedDir")
			return err
		}
		err = s.DirRepo.Update(tx, prevDirId, models.Directory{
			ParentDirId: &topDir.DirId,
			Name:        folder,
		})
		if err != nil {
			log.Warn().Err(err).Int("subDirId", prevDirId).Msg("Failed to update subdirectory")
			return err
		}
	}

	return nil
}

func (s *Service) isEqualDirTracked(tx *sqlx.Tx, absolutePath string) (equalDirTracked bool, err error) {
	trackedDirs, err := s.DirRepo.ReadAll(tx)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to fetch all dirs")
		return false, err
	}

	for _, trackedDir := range trackedDirs {
		trackedDirPath, err := s.AbsolutePath(tx, trackedDir.DirId)
		if err != nil {
			log.Warn().Err(err).Int("trackedDirId", trackedDir.DirId).Msg("Failed to calculate absolute path")
			return false, err
		}
		if absolutePath == trackedDirPath {
			return true, nil
		}
	}
	return false, nil
}

func (s *Service) getEqualTrackedDir(tx *sqlx.Tx, absolutePath string) (equalTrackedDir models.Directory, err error) {
	trackedDirs, err := s.DirRepo.ReadAll(tx)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to fetch all dirs")
		return models.Directory{}, err
	}

	for _, trackedDir := range trackedDirs {
		trackedDirPath, err := s.AbsolutePath(tx, trackedDir.DirId)
		if err != nil {
			log.Warn().Err(err).Int("trackedDirId", trackedDir.DirId).Msg("Failed to calculate absolute path")
			return models.Directory{}, err
		}
		if absolutePath == trackedDirPath {
			return trackedDir, nil
		}
	}
	return models.Directory{}, nil
}
