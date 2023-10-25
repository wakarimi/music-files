package song_service

import (
	"github.com/jmoiron/sqlx"
)

func (s *Service) IsExistsByDirAndName(tx *sqlx.Tx, dirId int, name string) (exists bool, err error) {
	exists, err = s.SongRepo.IsExistsByDirAndName(tx, dirId, name)
	if err != nil {
		return false, err
	}

	return exists, nil
}
