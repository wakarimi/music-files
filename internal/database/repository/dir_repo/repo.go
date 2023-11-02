package dir_repo

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/model"
)

type Repo interface {
	Create(tx *sqlx.Tx, dir model.Directory) (dirId int, err error)
	Read(tx *sqlx.Tx, dirId int) (dir model.Directory, err error)
	ReadAll(tx *sqlx.Tx) (dirs []model.Directory, err error)
	ReadRoots(tx *sqlx.Tx) (dirs []model.Directory, err error)
	ReadSubDirs(tx *sqlx.Tx, parentDirId int) (dirs []model.Directory, err error)
	ReadByParentAndName(tx *sqlx.Tx, parentDirId *int, name string) (dir model.Directory, err error)
	Update(tx *sqlx.Tx, dirId int, dir model.Directory) (err error)
	Delete(tx *sqlx.Tx, dirId int) (err error)
	IsExists(tx *sqlx.Tx, dirId int) (exists bool, err error)
	IsExistsByParentAndName(tx *sqlx.Tx, parentDirId *int, name string) (exists bool, err error)
}

type Repository struct {
}

func NewRepository() Repo {
	return &Repository{}
}
