package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/utils"
)

func (s *Service) DeleteOrphaned(tx *sqlx.Tx) (err error) {
	dirs, err := s.DirRepo.ReadAll(tx)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	for _, dir := range dirs {
		absolutePath, err := s.AbsolutePath(tx, dir.DirId)
		if err != nil {
			log.Error().Err(err)
			return err
		}

		existInDb, err := s.DirRepo.IsExists(tx, dir.ParentDirId, dir.Name)
		if err != nil {
			log.Error().Err(err)
			return err
		}

		if existInDb {
			existsOnDisk, err := utils.DirectoryExists(absolutePath)
			if err != nil {
				log.Error().Err(err)
				return err
			}

			if !existsOnDisk {
				err := s.Delete(tx, dir.DirId)
				if err != nil {
					log.Error().Err(err)
					return err
				}
			}
		}
	}

	return nil
}
