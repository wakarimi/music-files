package use_case

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/model/directory"
)

type transactor interface {
	WithTransaction(do func(tx *sqlx.Tx) (err error)) (err error)
}

type audioService interface {
	DeleteAllByDir(tx *sqlx.Tx, dirID int) error
}

type coverService interface {
	DeleteAllByDir(tx *sqlx.Tx, dirID int) error
}

type dirService interface {
	IsTracked(tx *sqlx.Tx, path string) (bool, error)
	IsExistsOnDisk(path string) (bool, error)
	Create(tx *sqlx.Tx, dirToCreate directory.Directory) (int, error)
	Get(tx *sqlx.Tx, dirID int) (directory.Directory, error)
	ContainedRoots(tx *sqlx.Tx, path string) ([]directory.Directory, error)
	MergeRoots(tx *sqlx.Tx, dirID1 int, dirID2 int) error
	IsExists(tx *sqlx.Tx, dirID int) (bool, error)
	IsRoot(tx *sqlx.Tx, dirID int) (bool, error)
	GetSubDirs(tx *sqlx.Tx, dirID int) ([]directory.Directory, error)
	Delete(tx *sqlx.Tx, dirID int) error
	GetRoots(tx *sqlx.Tx) ([]directory.Directory, error)
}

type UseCase struct {
	transactor   transactor
	audioService audioService
	coverService coverService
	dirService   dirService
}

func New(transactor transactor,
	audioService audioService,
	coverService coverService,
	dirService dirService) *UseCase {
	return &UseCase{
		transactor:   transactor,
		audioService: audioService,
		coverService: coverService,
		dirService:   dirService,
	}
}
