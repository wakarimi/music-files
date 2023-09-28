package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
	"music-files/internal/utils"
	"os"
)

func (s *Service) Scan(tx *sqlx.Tx, dirId int) (err error) {
	log.Debug().Int("dirId", dirId).
		Msg("Scanning directory")

	absolutePath, err := s.AbsolutePath(tx, dirId)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	exists, err := utils.DirectoryExists(absolutePath)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	if !exists {
		log.Info().Err(err).Str("path", absolutePath).
			Msg("Directory not exists")
		return nil
	}

	entries, err := os.ReadDir(absolutePath)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId)
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			alreadyInDb, err := s.DirRepo.IsExistsByParentAndName(tx, &dirId, entry.Name())
			if err != nil {
				log.Error().Err(err).Int("dirId", dirId)
				return err
			}

			var nextDirIdToScan int
			if alreadyInDb {
				nextDirToScan, err := s.DirRepo.ReadByParentAndName(tx, &dirId, entry.Name())
				if err != nil {
					log.Error().Err(err).Interface("dirId", dirId)
					return err
				}
				nextDirIdToScan = nextDirToScan.DirId
			} else {
				dirToCreate := models.Directory{
					ParentDirId: &dirId,
					Name:        entry.Name(),
				}
				nextDirIdToScan, err = s.DirRepo.Create(tx, dirToCreate)
			}

			err = s.Scan(tx, nextDirIdToScan)
			if err != nil {
				log.Error().Err(err)
				return err
			}
		}
	}

	log.Debug().Int("dirId", dirId).
		Msg("Directory scanned successfully")
	return nil
}
