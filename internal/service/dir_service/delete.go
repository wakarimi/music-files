package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s *Service) Delete(tx *sqlx.Tx, dirId int) (err error) {
	err = s.TrackRepo.DeleteByDirIdTx(tx, dirId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete tracks associated with the directory")
		return err
	}

	err = s.CoverRepo.DeleteByDirIdTx(tx, dirId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete covers associated with the directory")
		return err
	}

	err = s.DirRepo.DeleteTx(tx, dirId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete directory")
		return err
	}

	return nil
}
