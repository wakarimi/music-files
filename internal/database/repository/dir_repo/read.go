package dir_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (r *Repository) Read(tx *sqlx.Tx, dirId int) (dir models.Directory, err error) {
	log.Debug().Int("dirId", dirId).Msg("Fetching directory by id")

	query := `
		SELECT *
		FROM directories
		WHERE dir_id = :dir_id
	`
	args := map[string]interface{}{
		"dir_id": dirId,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId)
		return models.Directory{}, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close rows")
		}
	}(rows)
	if rows.Next() {
		if err = rows.StructScan(&dir); err != nil {
			log.Error().Err(err).Int("dirId", dirId)
			return models.Directory{}, err
		}
	} else {
		err := fmt.Errorf("no directory found with dir_id: %d", dirId)
		log.Error().Err(err)
		return models.Directory{}, err
	}

	log.Debug().Str("name", dir.Name).Msg("Directory fetched by id successfully")
	return dir, nil
}
