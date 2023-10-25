package cover_repo

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/models"
)

type Repo interface {
	Create(tx *sqlx.Tx, cover models.Cover) (coverId int, err error)
	Read(tx *sqlx.Tx, coverId int) (cover models.Cover, err error)
	ReadByDirAndName(tx *sqlx.Tx, dirId int, name string) (cover models.Cover, err error)
	ReadAllByDir(tx *sqlx.Tx, dirId int) (covers []models.Cover, err error)
	Update(tx *sqlx.Tx, coverId int, cover models.Cover) (err error)
	Delete(tx *sqlx.Tx, coverId int) (err error)
	IsExists(tx *sqlx.Tx, coverId int) (exists bool, err error)
	IsExistsByDirAndName(tx *sqlx.Tx, dirId int, name string) (exists bool, err error)
}

type Repository struct {
}

func NewRepository() Repo {
	return &Repository{}
}
