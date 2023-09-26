package dir_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r *Repository) Delete(tx *sqlx.Tx, dirId int) (err error) {
	log.Debug().Int("dirId", dirId).Msg("Deleting directory")

	query := `
		DELETE FROM directories
		WHERE dir_id = :dir_id
	`
	args := map[string]interface{}{
		"dir_id": dirId,
	}
	_, err = tx.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Msg("Failed to delete directory")
		return err
	}

	log.Debug().Int("dirId", dirId).Msg("Directory deleted successfully")
	return nil
}
