package dir_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (r *Repository) ReadSubDirs(tx *sqlx.Tx, parentDirId int) (dirs []models.Directory, err error) {
	log.Debug().Int("parentDirId", parentDirId).Msg("Fetching subdirectories")

	query := `
		SELECT * 
		FROM directories
		WHERE parent_dir_id = :parentDirId
	`
	args := map[string]interface{}{
		"parentDirId": parentDirId,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err)
		}
	}(rows)

	log.Debug().Int("parentDirId", parentDirId).Int("subDirsCount", len(dirs)).Msg("Subdirectories fetched successfully")
	return dirs, nil
}
