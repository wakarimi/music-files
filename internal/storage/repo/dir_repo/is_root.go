package dir_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) IsRoot(tx *sqlx.Tx, dirID int) (bool, error) {
	log.Debug().Int("dirId", dirID).Msg("Checking if directory is root")

	query := `
		SELECT NOT EXISTS (
			SELECT 1 
			FROM directories
			WHERE id = $1
			AND parent_dir_id IS NOT NULL
		)
	`

	var isRoot bool
	err := tx.Get(&isRoot, query, dirID)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirID).Msg("Failed to execute query to check if directory is root")
		return false, err
	}

	log.Debug().Int("dirId", dirID).Bool("isRoot", isRoot).Msg("Directory root check successful")
	return isRoot, nil
}
