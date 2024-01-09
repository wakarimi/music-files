package audio_service

import (
	"github.com/jmoiron/sqlx"
)

type audioRepo interface {
	DeleteByDir(tx *sqlx.Tx, dirID int) error
}

type Service struct {
	audioRepo audioRepo
}

func New(audioRepo audioRepo) *Service {
	return &Service{
		audioRepo: audioRepo,
	}
}
