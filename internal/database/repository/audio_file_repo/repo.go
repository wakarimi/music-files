package audio_file_repo

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/model"
)

type Repo interface {
	Create(tx *sqlx.Tx, audioFile model.AudioFile) (audioFileId int, err error)
	Read(tx *sqlx.Tx, audioFileId int) (audioFile model.AudioFile, err error)
	ReadByDirAndName(tx *sqlx.Tx, dirId int, name string) (audioFile model.AudioFile, err error)
	ReadAll(tx *sqlx.Tx) (audioFiles []model.AudioFile, err error)
	ReadAllBySha256(tx *sqlx.Tx, sha256 string) (audioFiles []model.AudioFile, err error)
	ReadAllByDir(tx *sqlx.Tx, dirId int) (audioFiles []model.AudioFile, err error)
	Update(tx *sqlx.Tx, audioFileId int, audioFile model.AudioFile) (err error)
	Delete(tx *sqlx.Tx, audioFileId int) (err error)
	IsExists(tx *sqlx.Tx, audioFileId int) (exists bool, err error)
	IsExistsByDirAndName(tx *sqlx.Tx, dirId int, name string) (exists bool, err error)
}

type Repository struct {
}

func NewRepository() Repo {
	return &Repository{}
}
