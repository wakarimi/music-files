package audio_service

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/model/audio"
)

type audioRepo interface {
	DeleteByDir(tx *sqlx.Tx, dirID int) error
	Create(tx *sqlx.Tx, audioToCreate audio.Audio) (int, error)
	Delete(tx *sqlx.Tx, audioID int) error
	ReadAllByDir(tx *sqlx.Tx, dirID int) ([]audio.Audio, error)
	ReadByDirAndName(tx *sqlx.Tx, dirID int, name string) (audio.Audio, error)
	IsExistsByDirAndName(tx *sqlx.Tx, dirID int, name string) (bool, error)
	Update(tx *sqlx.Tx, audioID int, update audio.Audio) error
	Read(tx *sqlx.Tx, audioID int) (audio.Audio, error)
	IsExists(tx *sqlx.Tx, audioID int) (bool, error)
}

type Service struct {
	audioRepo audioRepo
}

func New(audioRepo audioRepo) *Service {
	return &Service{
		audioRepo: audioRepo,
	}
}
