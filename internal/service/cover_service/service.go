package cover_service

import (
	"github.com/jmoiron/sqlx"
)

type coverRepo interface {
	DeleteByDir(tx *sqlx.Tx, dirID int) error
}

type Service struct {
	coverRepo coverRepo
}

func New(coverRepo coverRepo) *Service {
	return &Service{
		coverRepo: coverRepo,
	}
}
