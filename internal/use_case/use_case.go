package use_case

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/model/audio"
	"music-files/internal/model/cover"
	"music-files/internal/model/directory"
)

type transactor interface {
	WithTransaction(do func(tx *sqlx.Tx) (err error)) (err error)
}

type audioService interface {
	DeleteAllByDir(tx *sqlx.Tx, dirID int) error
	IsAudioByPath(path string) (bool, error)
	CalculateSHA256(path string) (string, error)
	IsExistsByDirAndName(tx *sqlx.Tx, dirID int, name string) (bool, error)
	GetByDirAndName(tx *sqlx.Tx, dirID int, name string) (audio.Audio, error)
	ConstructByPath(path string) (audio.Audio, error)
	Update(tx *sqlx.Tx, audioID int, audioToUpdate audio.Audio) error
	Create(tx *sqlx.Tx, audioToCreate audio.Audio) (int, error)
	GetAllByDir(tx *sqlx.Tx, dirID int) ([]audio.Audio, error)
	Delete(tx *sqlx.Tx, audioID int) error
}

type coverService interface {
	DeleteAllByDir(tx *sqlx.Tx, dirID int) error
	IsCoverByPath(path string) (bool, error)
	CalculateSHA256(path string) (string, error)
	IsExistsByDirAndName(tx *sqlx.Tx, dirID int, name string) (bool, error)
	GetByDirAndName(tx *sqlx.Tx, dirID int, name string) (cover.Cover, error)
	ConstructByPath(path string) (cover.Cover, error)
	Update(tx *sqlx.Tx, coverID int, coverToUpdate cover.Cover) error
	Create(tx *sqlx.Tx, coverToCreate cover.Cover) (int, error)
	GetAllByDir(tx *sqlx.Tx, dirID int) ([]cover.Cover, error)
	Delete(tx *sqlx.Tx, coverID int) error
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
	CalcAbsolutePath(tx *sqlx.Tx, dirID int) (string, error)
	IsExistsByParentAndName(tx *sqlx.Tx, parentID *int, name string) (bool, error)
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
