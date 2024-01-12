package cover_service

import "github.com/jmoiron/sqlx"

func (s Service) IsExistsByDirAndName(tx *sqlx.Tx, dirID int, name string) (bool, error) {
	//TODO implement me
	panic("implement me")
}
