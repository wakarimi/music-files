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
		WHERE parent_dir_id = :parent_dir_id
	`
	args := map[string]interface{}{
		"parent_dir_id": parentDirId,
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

	for rows.Next() {
		var dir models.Directory
		if err = rows.StructScan(&dir); err != nil {
			log.Error().Err(err)
			return nil, err
		}
		dirs = append(dirs, dir)
	}

	log.Debug().Int("parentDirId", parentDirId).Int("subDirsCount", len(dirs)).Msg("Subdirectories fetched successfully")
	return dirs, nil
}
