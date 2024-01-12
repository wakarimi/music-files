package cover_service

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/model/cover"
)

func (s Service) Update(tx *sqlx.Tx, coverID int, coverToUpdate cover.Cover) error {
	//TODO implement me
	panic("implement me")
}
