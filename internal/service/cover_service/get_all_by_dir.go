package cover_service

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/models"
)

func (s *Service) GetAllByDir(tx *sqlx.Tx, dirId int) (covers []models.Cover, err error) {
	covers, err = s.CoverRepo.ReadAllByDir(tx, dirId)
	if err != nil {
		return make([]models.Cover, 0), err
	}

	return covers, nil
}
