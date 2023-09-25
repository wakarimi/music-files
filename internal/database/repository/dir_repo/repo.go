package dir_repo

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/models"
)

type Repo interface {
	Create(tx *sqlx.Tx, dir models.Directory) (dirId int, err error)
	IsExists(tx *sqlx.Tx, parentDirId int, name string) (exists bool, err error)
}

type Repository struct {
}

func NewRepository() Repo {
	return &Repository{}
}
