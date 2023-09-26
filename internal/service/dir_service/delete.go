package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s *Service) Delete(tx *sqlx.Tx, dirId int) (err error) {
	subDirs, err := s.DirRepo.ReadSubDirs(tx, dirId)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	for _, subDir := range subDirs {
		err := s.Delete(tx, subDir.DirId)
		if err != nil {
			log.Error().Err(err)
			return err
		}
	}

	err = s.DirRepo.Delete(tx, dirId)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	return nil
}
