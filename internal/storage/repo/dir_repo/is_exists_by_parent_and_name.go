package dir_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) IsExistsByParentAndName(tx *sqlx.Tx, parentDirID *int, name string) (exists bool, err error) {
	log.Debug().Interface("parentDirId", parentDirID).Str("name", name).Msg("Checking if directory exists in database")

	var query string
	var row *sqlx.Rows
	if parentDirID == nil {
		query = `
            SELECT EXISTS (
                SELECT 1 
                FROM directories
                WHERE parent_dir_id IS NULL
                	AND name = $1
            )
        `
		row, err = tx.Queryx(query, name)
	} else {
		query = `
            SELECT EXISTS (
                SELECT 1 
                FROM directories
                WHERE parent_dir_id = $1
                	AND name = $2
            )
        `
		row, err = tx.Queryx(query, *parentDirID, name)
	}
	if err != nil {
		log.Error().Err(err).Interface("parentDirId", parentDirID).Str("name", name)
		return false, err
	}
	defer func(row *sqlx.Rows) {
		err := row.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close row")
		}
	}(row)

	if row.Next() {
		if err = row.Scan(&exists); err != nil {
			log.Error().Err(err).Interface("parentDirId", parentDirID).Str("name", name)
			return false, err
		}
	}

	log.Debug().Interface("parentDirId", parentDirID).Str("name", name).Bool("exists", exists).Msg("The existence of the directory was checked successfully")
	return exists, nil
}
