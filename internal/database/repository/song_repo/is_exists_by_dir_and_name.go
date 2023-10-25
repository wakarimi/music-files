package song_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) IsExistsByDirAndName(tx *sqlx.Tx, dirId int, name string) (exists bool, err error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM songs
			WHERE dir_id = :dir_id
				AND filename = :name
		)
	`
	args := map[string]interface{}{
		"dir_id": dirId,
		"name":   name,
	}
	row, err := tx.NamedQuery(query, args)
	if err != nil {
		return false, err
	}
	defer func(row *sqlx.Rows) {
		err := row.Close()
		if err != nil {
			log.Error().Err(err)
		}
	}(row)
	if row.Next() {
		if err = row.Scan(&exists); err != nil {
			return false, err
		}
	}

	return exists, nil
}
