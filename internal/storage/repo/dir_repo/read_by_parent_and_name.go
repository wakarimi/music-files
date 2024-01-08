package dir_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/directory"
)

func (r Repository) ReadByParentAndName(tx *sqlx.Tx, parentDirID *int, name string) (dir directory.Directory, err error) {
	log.Debug().Interface("parentDirId", parentDirID).Str("name", name).Msg("Fetching directory by parent and name")

	var query string
	var row *sqlx.Rows
	if parentDirID == nil {
		query = `
		SELECT *
		FROM directories
		WHERE parent_dir_id IS NULL
			AND name = $1
    `
		row, err = tx.Queryx(query, name)
	} else {
		query = `
		SELECT *
		FROM directories
		WHERE parent_dir_id = :parent_dir_id
			AND name = :name
    `
		args := map[string]interface{}{
			"parent_dir_id": *parentDirID,
			"name":          name,
		}
		row, err = tx.NamedQuery(query, args)
	}
	if err != nil {
		log.Error().Err(err).Interface("parentDirId", parentDirID).Str("name", name)
		return directory.Directory{}, err
	}
	defer func(row *sqlx.Rows) {
		err := row.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close row")
		}
	}(row)

	if row.Next() {
		if err = row.StructScan(&dir); err != nil {
			log.Error().Interface("parentDirId", parentDirID).Str("name", name)
			return directory.Directory{}, err
		}
	}

	log.Debug().Interface("parentDirId", parentDirID).Str("name", name).Msg("Directory fetched by id successfully")
	return dir, nil
}
