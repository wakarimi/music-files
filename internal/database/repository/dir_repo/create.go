package dir_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (r *Repository) Create(tx *sqlx.Tx, dir models.Directory) (dirId int, err error) {
	log.Debug().Interface("parentDirId", dir.ParentDirId).Str("name", dir.Name).Msg("Creating new directory")

	query := `
		INSERT INTO directories(name, parent_dir_id)
		VALUES (:name, :parent_dir_id)
		RETURNING dir_id
	`
	rows, err := tx.NamedQuery(query, dir)
	if err != nil {
		log.Error().Err(err).Interface("parentDirId", dir.ParentDirId).Str("name", dir.Name)
		return 0, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close rows")
		}
	}(rows)

	if rows.Next() {
		if err := rows.Scan(&dirId); err != nil {
			log.Error().Err(err).Interface("parentDirId", dir.ParentDirId).Str("name", dir.Name)
			return 0, err
		}
	} else {
		err := fmt.Errorf("no id returned after directory insert")
		log.Error().Err(err).Interface("parentDirId", dir.ParentDirId).Str("name", dir.Name)
		return 0, err
	}

	log.Debug().Interface("parentDirId", dir.ParentDirId).Str("name", dir.Name).Int("dirId", dirId).Msg("Directory created successfully")
	return dirId, nil
}
