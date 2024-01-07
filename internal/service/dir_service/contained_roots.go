package dir_service

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/model/directory"
)

func (s Service) ContainedRoots(tx *sqlx.Tx, path string) ([]directory.Directory, error) {
	//TODO implement me
	panic("implement me")
}
