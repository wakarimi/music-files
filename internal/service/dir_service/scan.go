package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
	"music-files/internal/utils"
	"os"
	"path/filepath"
)

func (s *Service) Scan(tx *sqlx.Tx, dirId int) (err error) {
	log.Debug().Int("dirId", dirId).
		Msg("Scanning directory")

	dirAbsolutePath, err := s.AbsolutePath(tx, dirId)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	exists, err := utils.DirectoryExists(dirAbsolutePath)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	if !exists {
		log.Info().Err(err).Str("path", dirAbsolutePath).
			Msg("Directory not exists")
		return nil
	}

	entries, err := os.ReadDir(dirAbsolutePath)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId)
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			err = s.actualizeDir(tx, dirId, entry.Name())
			if err != nil {
				return err
			}
			continue
		}

		fileAbsolutePath := filepath.Join(dirAbsolutePath, entry.Name())
		isMusicFile, err := utils.IsMusicFile(fileAbsolutePath)
		if err != nil {
			return err
		}
		isImageFile, err := utils.IsImageFile(fileAbsolutePath)
		if err != nil {
			return err
		}

		if isMusicFile {
			err = s.actualizeMusicFile(tx, dirId, entry.Name())
			if err != nil {
				return err
			}
		} else if isImageFile {
			err = s.actualizeImageFile(tx, dirId, entry.Name())
			if err != nil {
				return err
			}
		}
	}

	log.Debug().Int("dirId", dirId).
		Msg("Directory scanned successfully")
	return nil
}

func (s *Service) actualizeDir(tx *sqlx.Tx, dirId int, name string) error {
	alreadyInDb, err := s.DirRepo.IsExistsByParentAndName(tx, &dirId, name)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId)
		return err
	}

	var nextDirIdToScan int
	if alreadyInDb {
		nextDirToScan, err := s.DirRepo.ReadByParentAndName(tx, &dirId, name)
		if err != nil {
			log.Error().Err(err).Interface("dirId", dirId)
			return err
		}
		nextDirIdToScan = nextDirToScan.DirId
	} else {
		dirToCreate := models.Directory{
			ParentDirId: &dirId,
			Name:        name,
		}
		nextDirIdToScan, err = s.DirRepo.Create(tx, dirToCreate)
	}

	err = s.Scan(tx, nextDirIdToScan)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	return nil
}

func (s *Service) actualizeMusicFile(tx *sqlx.Tx, dirId int, name string) error {
	alreadyInDb, err := s.DirRepo.IsExistsByParentAndName(tx, &dirId, name)

	return nil
}
func (s *Service) actualizeImageFile(tx *sqlx.Tx, dirId int, name string) error {
	alreadyInDb, err := s.DirRepo.IsExistsByParentAndName(tx, &dirId, name)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId)
		return err
	}

	return nil
}
