package dir_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/directory"
)

func (r Repository) Update(tx *sqlx.Tx, dirID int, dir directory.Directory) (err error) {
	log.Debug().Int("dirId", dirID).Interface("parentDirId", dir.ParentDirID).Str("name", dir.Name).Msg("Updating directory")

	query := `
		UPDATE directories
		SET name = :name, parent_dir_id = :parent_dir_id
		WHERE id = :id
	`

	dir.ID = dirID
	_, err = tx.NamedExec(query, dir)

	if err != nil {
		log.Error().Err(err).Int("dirId", dirID).Interface("parentDirId", dir.ParentDirID).Str("name", dir.Name)
		return err
	}

	log.Debug().Int("dirId", dirID).Interface("parentDirId", dir.ParentDirID).Str("name", dir.Name).Msg("Directory updated successfully")
	return nil
}
