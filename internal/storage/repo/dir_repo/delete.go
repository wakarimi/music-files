package dir_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) Delete(tx *sqlx.Tx, dirID int) (err error) {
	log.Debug().Int("dirId", dirID).Msg("Deleting directory")

	query := `
		DELETE FROM directories
		WHERE id = :id
	`
	args := map[string]interface{}{
		"id": dirID,
	}
	_, err = tx.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirID).Msg("Failed to delete directory")
		return err
	}

	log.Debug().Int("dirId", dirID).Msg("Directory deleted successfully")
	return nil
}
