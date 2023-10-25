package song_repo

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/models"
)

type Repo interface {
	Create(tx *sqlx.Tx, song models.Song) (songId int, err error)
	Read(tx *sqlx.Tx, songId int) (song models.Song, err error)
	ReadByDirAndName(tx *sqlx.Tx, dirId int, name string) (song models.Song, err error)
	ReadAllByDir(tx *sqlx.Tx, dirId int) (songs []models.Song, err error)
	Update(tx *sqlx.Tx, songId int, song models.Song) (err error)
	Delete(tx *sqlx.Tx, songId int) (err error)
	IsExists(tx *sqlx.Tx, songId int) (exists bool, err error)
	IsExistsByDirAndName(tx *sqlx.Tx, dirId int, name string) (exists bool, err error)
}

type Repository struct {
}

func NewRepository() Repo {
	return &Repository{}
}
