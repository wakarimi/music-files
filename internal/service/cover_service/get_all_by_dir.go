package cover_service

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/model/cover"
)

func (s Service) GetAllByDir(tx *sqlx.Tx, dirID int) ([]cover.Cover, error) {
	//TODO implement me
	panic("implement me")
}
