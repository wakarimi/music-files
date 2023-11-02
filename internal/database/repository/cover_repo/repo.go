package cover_repo

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/model"
)

type Repo interface {
	Create(tx *sqlx.Tx, cover model.Cover) (coverId int, err error)
	Read(tx *sqlx.Tx, coverId int) (cover model.Cover, err error)
	ReadByDirAndName(tx *sqlx.Tx, dirId int, name string) (cover model.Cover, err error)
	ReadAllByDir(tx *sqlx.Tx, dirId int) (covers []model.Cover, err error)
	Update(tx *sqlx.Tx, coverId int, cover model.Cover) (err error)
	Delete(tx *sqlx.Tx, coverId int) (err error)
	IsExists(tx *sqlx.Tx, coverId int) (exists bool, err error)
	IsExistsByDirAndName(tx *sqlx.Tx, dirId int, name string) (exists bool, err error)
}

type Repository struct {
}

func NewRepository() Repo {
	return &Repository{}
}
