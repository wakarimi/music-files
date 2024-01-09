package dir_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/directory"
)

func (r Repository) ReadSubDirs(tx *sqlx.Tx, dirID int) (dirs []directory.Directory, err error) {
	log.Debug().Int("parentDirId", dirID).Msg("Fetching subdirectories")

	query := `
		SELECT * 
		FROM directories
		WHERE parent_dir_id = :parent_dir_id
	`
	args := map[string]interface{}{
		"parent_dir_id": dirID,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("parentDirId", dirID).Str("query", query).Msg("Failed to execute query to read subdirectories")
		return nil, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close rows")
		}
	}(rows)

	for rows.Next() {
		var dir directory.Directory
		if err = rows.StructScan(&dir); err != nil {
			log.Error().Err(err).Int("parentDirId", dirID).Msg("Failed to get read result")
			return nil, err
		}
		dirs = append(dirs, dir)
	}

	log.Debug().Int("parentDirId", dirID).Int("subDirsCount", len(dirs)).Msg("Subdirectories fetched successfully")
	return dirs, nil
}
