package dir_service

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/model/directory"
)

type dirRepo interface {
	ReadRoots(tx *sqlx.Tx) ([]directory.Directory, error)
	Read(tx *sqlx.Tx, dirID int) (directory.Directory, error)
	Create(tx *sqlx.Tx, dirToCreate directory.Directory) (int, error)
	IsExistsByParentAndName(tx *sqlx.Tx, parentID *int, name string) (bool, error)
	ReadByParentAndName(tx *sqlx.Tx, parentID *int, name string) (directory.Directory, error)
	Update(tx *sqlx.Tx, dirID int, directory directory.Directory) error
}

type Service struct {
	dirRepo dirRepo
}

func New(dirRepo dirRepo) *Service {
	return &Service{
		dirRepo: dirRepo,
	}
}
