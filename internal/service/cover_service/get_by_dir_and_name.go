package cover_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"music-files/internal/errors"
	"music-files/internal/models"
)

func (s *Service) GetByDirAndName(tx *sqlx.Tx, dirId int, name string) (cover models.Cover, err error) {
	exists, err := s.CoverRepo.IsExistsByDirAndName(tx, dirId, name)
	if err != nil {
		return models.Cover{}, err
	}
	if !exists {
		return models.Cover{}, errors.NotFound{Resource: fmt.Sprintf("cover with dirId=%d and name=%s in database", dirId, name)}
	}

	cover, err = s.CoverRepo.ReadByDirAndName(tx, dirId, name)
	if err != nil {
		return models.Cover{}, err
	}

	return cover, nil
}
