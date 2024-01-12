package cover_service

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/model/cover"
)

func (s Service) Create(tx *sqlx.Tx, coverToCreate cover.Cover) (int, error) {
	//TODO implement me
	panic("implement me")
}
