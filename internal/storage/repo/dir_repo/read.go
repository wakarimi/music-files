package dir_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/directory"
)

func (r Repository) Read(tx *sqlx.Tx, dirID int) (dir directory.Directory, err error) {
	log.Debug().Int("dirId", dirID).Msg("Fetching directory by id")

	query := `
		SELECT *
		FROM directories
		WHERE id = :id
	`
	args := map[string]interface{}{
		"id": dirID,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirID)
		return directory.Directory{}, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close rows")
		}
	}(rows)
	if rows.Next() {
		if err = rows.StructScan(&dir); err != nil {
			log.Error().Err(err).Int("dirId", dirID)
			return directory.Directory{}, err
		}
	} else {
		err := fmt.Errorf("no directory found with dir_id: %d", dirID)
		log.Error().Err(err)
		return directory.Directory{}, err
	}

	log.Debug().Str("name", dir.Name).Msg("Directory fetched by id successfully")
	return dir, nil
}
