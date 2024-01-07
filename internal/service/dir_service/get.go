package dir_service

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/model/directory"
)

func (s Service) Get(tx *sqlx.Tx, dirID int) (directory.Directory, error) {
	//TODO implement me
	panic("implement me")
}
