package dir_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (r *Repository) ReadAll(tx *sqlx.Tx) (dirs []models.Directory, err error) {
	log.Debug().Msg("Fetching all directories")

	query := `
		SELECT * 
		FROM directories
	`
	err = tx.Select(&dirs, query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch directories")
		return nil, err
	}

	log.Debug().Int("dirsCount", len(dirs)).Msg("All directories fetched successfully")
	return dirs, nil
}
