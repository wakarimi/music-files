package cover_service

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/model/cover"
)

func (s Service) GetByDirAndName(tx *sqlx.Tx, dirID int, name string) (cover.Cover, error) {
	//TODO implement me
	panic("implement me")
}
