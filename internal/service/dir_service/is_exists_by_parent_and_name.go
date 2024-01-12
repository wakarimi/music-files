package dir_service

import "github.com/jmoiron/sqlx"

func (s Service) IsExistsByParentAndName(tx *sqlx.Tx, parentID *int, name string) (bool, error) {
	//TODO implement me
	panic("implement me")
}
