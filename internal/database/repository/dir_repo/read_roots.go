package dir_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (r *Repository) ReadRoots(tx *sqlx.Tx) (dirs []models.Directory, err error) {
	log.Debug().Msg("Reading root directories")

	query := `
		SELECT * 
		FROM directories
		WHERE parent_dir_id IS NULL
	`
	err = tx.Select(&dirs, query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read root directories")
		return nil, err
	}

	log.Debug().Int("countOfRootDirs", len(dirs)).Msg("All root directories read successfully")
	return dirs, nil
}
