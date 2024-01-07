package dir_service

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/model/directory"
)

func (s Service) Create(tx *sqlx.Tx, dirToCreate directory.Directory) (int, error) {
	//TODO implement me
	panic("implement me")
}
