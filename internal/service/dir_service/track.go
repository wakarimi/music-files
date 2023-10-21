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

	_, err = os.Stat(dir.Name)
	if err != nil {
		log.Error().Err(err).Msg("Directory on disk not found")
	}
	if os.IsNotExist(err) {
		return models.Directory{}, errors.NotFound{Resource: "directory on disk"}
	} else if err != nil {
		return models.Directory{}, err
	}

	alreadyTracked, err := s.isAlreadyTracked(tx, dir.Name)
	if err != nil {
		log.Error().Err(err).Str("path", dir.Name).Msg("Failed to check whether directory was added")
		return models.Directory{}, err
	}
	if alreadyTracked {
		err = fmt.Errorf("directory with path=%s already tracked", dir.Name)
		log.Error().Err(err).Str("path", dir.Name).Msg("Directory already tracked")
		return models.Directory{}, errors.Conflict{Message: err.Error()}
	}

	createdDirId, err := s.DirRepo.Create(tx, dir)
	if err != nil {
		log.Error().Err(err).Str("path", dir.Name).Msg("Failed to track directory")
		return models.Directory{}, err
	}
	createdDir, err = s.DirRepo.Read(tx, createdDirId)
	if err != nil {
		log.Error().Err(err).Int("dirId", createdDir.DirId).Msg("Failed to read added directory")
		return models.Directory{}, err
	}

	containedRoots, err := s.getContainedRoots(tx, dir.Name)
	if err != nil {
		log.Error().Err(err).Str("path", dir.Name).Msg("Failed to get contained roots")
		return models.Directory{}, err
	}
	log.Debug().Int("containedRootsCount", len(containedRoots)).Msg("containedRootsCount:")
	err = s.addIntermediateDirs(tx, createdDir, containedRoots)
	if err != nil {
		log.Error().Err(err).Str("path", dir.Name).Msg("Failed to add intermediate dirs")
		return models.Directory{}, err
	}

	log.Debug().Str("path", dir.Name).Msg("Directory added to tracked")
	return createdDir, nil
}

func (s *Service) isAlreadyTracked(tx *sqlx.Tx, absolutePath string) (alreadyTracked bool, err error) {
	trackedDirs, err := s.DirRepo.ReadAll(tx)
	if err != nil {
		return false, err
	}

	for _, trackedDir := range trackedDirs {
		trackedDirPath, err := s.AbsolutePath(tx, trackedDir.DirId)
		if err != nil {
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
		return make([]models.Directory, 0), err
	}

	log.Debug().Str("absolutePath", absolutePath).Msg("delme")
	for _, root := range roots {
		if strings.Contains(root.Name, absolutePath) && root.Name != absolutePath {
			containedRoots = append(containedRoots, root)
		}
	}
	return containedRoots, nil
}

func (s *Service) addIntermediateDirs(tx *sqlx.Tx, root models.Directory, containedRoots []models.Directory) (err error) {
	for _, containedRoot := range containedRoots {
		prevDirId := containedRoot.DirId
		currentPath, folder := filepath.Split(containedRoot.Name)
		currentPath = filepath.Clean(currentPath)
		alreadyTracked, err := s.isEqualDirTracked(tx, currentPath)
		if err != nil {
			return err
		}
		for !alreadyTracked {
			newDirId, err := s.DirRepo.Create(tx, models.Directory{
				ParentDirId: nil,
				Name:        currentPath,
			})
			if err != nil {
				return err
			}
			err = s.DirRepo.Update(tx, prevDirId, models.Directory{
				ParentDirId: &newDirId,
				Name:        folder,
			})
			if err != nil {
				return err
			}
			prevDirId = newDirId
			currentPath, folder = filepath.Split(currentPath)
			currentPath = filepath.Clean(currentPath)
			alreadyTracked, err = s.isEqualDirTracked(tx, currentPath)
		}
		topDir, err := s.getEqualTrackedDir(tx, currentPath)
		if err != nil {
			return err
		}
		err = s.DirRepo.Update(tx, prevDirId, models.Directory{
			ParentDirId: &topDir.DirId,
			Name:        folder,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) isEqualDirTracked(tx *sqlx.Tx, absolutePath string) (equalDirTracked bool, err error) {
	trackedDirs, err := s.DirRepo.ReadAll(tx)
	if err != nil {
		return false, err
	}

	for _, trackedDir := range trackedDirs {
		trackedDirPath, err := s.AbsolutePath(tx, trackedDir.DirId)
		if err != nil {
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
		return models.Directory{}, err
	}

	for _, trackedDir := range trackedDirs {
		trackedDirPath, err := s.AbsolutePath(tx, trackedDir.DirId)
		if err != nil {
			return models.Directory{}, err
		}
		if absolutePath == trackedDirPath {
			return trackedDir, nil
		}
	}
	return models.Directory{}, nil
}
