package dir_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (r *Repository) Update(tx *sqlx.Tx, dirId int, dir models.Directory) (err error) {
	log.Debug().Int("dirId", dirId).Interface("parentDirId", dir.ParentDirId).Str("name", dir.Name).
		Msg("Updating directory")

	query := `
		UPDATE directories
		SET name = :name, parent_dir_id = :parent_dir_id
		WHERE dir_id = :dir_id
	`

	dir.DirId = dirId
	_, err = tx.NamedExec(query, dir)

	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Interface("parentDirId", dir.ParentDirId).Str("name", dir.Name)
		return err
	}

	log.Debug().Int("dirId", dirId).Interface("parentDirId", dir.ParentDirId).Str("name", dir.Name).
		Msg("Directory updated successfully")

	return nil
}
