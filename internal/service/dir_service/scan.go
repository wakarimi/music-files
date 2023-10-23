package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/models"
	"music-files/internal/utils"
	"os"
)

func (s *Service) Scan(tx *sqlx.Tx, dirId int) (err error) {
	log.Debug().Int("dirId", dirId).Msg("Scanning directory")

	existsInDatabase, err := s.DirRepo.IsExists(tx, dirId)
	if err != nil {
		return err
	}
	if !existsInDatabase {
		return errors.NotFound{Resource: "directory in database"}
	}

	absolutePath, err := s.AbsolutePath(tx, dirId)
	if err != nil {
		return err
	}
	existsOnDisk, err := utils.IsDirectoryExistsOnDisk(absolutePath)
	if err != nil {
		return err
	}
	if !existsOnDisk {
		err = s.DeleteDir(tx, dirId)
		if err != nil {
			return err
		}
		return nil
	}

	err = s.actualizeSubDirs(tx, dirId)
	if err != nil {
		return err
	}

	subDirs, err := s.DirRepo.ReadSubDirs(tx, dirId)
	if err != nil {
		return err
	}

	for _, subDir := range subDirs {
		err = s.Scan(tx, subDir.DirId)
		if err != nil {
			return err
		}
	}
	err = s.scanContent(tx, dirId)
	if err != nil {
		return err
	}

	log.Debug().Int("dirId", dirId).Msg("Directory scanned successfully")
	return nil
}

func (s *Service) actualizeSubDirs(tx *sqlx.Tx, dirId int) (err error) {
	absolutePath, err := s.AbsolutePath(tx, dirId)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(absolutePath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			alreadyInDatabase, err := s.DirRepo.IsExistsByParentAndName(tx, &dirId, entry.Name())
			if err != nil {
				return err
			}
			if !alreadyInDatabase {
				_, err = s.DirRepo.Create(tx, models.Directory{
					ParentDirId: &dirId,
					Name:        entry.Name(),
				})
				if err != nil {
					return err
				}
			}
		}
	}

	subDirs, err := s.DirRepo.ReadSubDirs(tx, dirId)
	if err != nil {
		return err
	}
	for _, subDir := range subDirs {
		foundDirOnDisk := false

		for _, entry := range entries {
			if entry.IsDir() {
				if subDir.Name == entry.Name() {
					foundDirOnDisk = true
				}
			}
		}

		if !foundDirOnDisk {
			err = s.DeleteDir(tx, subDir.DirId)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Service) scanContent(tx *sqlx.Tx, dirId int) (err error) {
	err = s.actualizeSongs(tx, dirId)
	if err != nil {
		return err
	}

	err = s.actualizeCovers(tx, dirId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) actualizeSongs(tx *sqlx.Tx, dirId int) (err error) {
	return nil
}

func (s *Service) actualizeCovers(tx *sqlx.Tx, dirId int) (err error) {
	return nil
}
