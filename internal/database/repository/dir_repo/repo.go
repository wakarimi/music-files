package dir_repo

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/models"
)

type Repo interface {
	Create(tx *sqlx.Tx, dir models.Directory) (dirId int, err error)
	Read(tx *sqlx.Tx, dirId int) (dir models.Directory, err error)
	ReadAll(tx *sqlx.Tx) (dirs []models.Directory, err error)
	ReadByParentAndName(tx *sqlx.Tx, parentDirId *int, name string) (dir models.Directory, err error)
	ReadSubDirs(tx *sqlx.Tx, parentDirId int) (dirs []models.Directory, err error)
	Delete(tx *sqlx.Tx, dirId int) (err error)
	IsExists(tx *sqlx.Tx, parentDirId *int, name string) (exists bool, err error)
}

type Repository struct {
}

func NewRepository() Repo {
	return &Repository{}
}
