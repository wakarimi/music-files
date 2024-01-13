package cover_service

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/model/cover"
)

type coverRepo interface {
	DeleteByDir(tx *sqlx.Tx, dirID int) error
	Create(tx *sqlx.Tx, coverToCreate cover.Cover) (int, error)
	Delete(tx *sqlx.Tx, coverID int) error
	ReadAllByDir(tx *sqlx.Tx, dirID int) ([]cover.Cover, error)
	ReadByDirAndName(tx *sqlx.Tx, dirID int, name string) (cover.Cover, error)
	IsExistsByDirAndName(tx *sqlx.Tx, dirID int, name string) (bool, error)
	Update(tx *sqlx.Tx, coverID int, coverToUpdate cover.Cover) error
}

type Service struct {
	coverRepo coverRepo
}

func New(coverRepo coverRepo) *Service {
	return &Service{
		coverRepo: coverRepo,
	}
}
