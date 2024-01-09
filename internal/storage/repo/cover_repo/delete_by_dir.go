package cover_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) DeleteByDir(tx *sqlx.Tx, dirID int) (err error) {
	log.Debug().Int("dirId", dirID).Msg("Deleting covers by directory")

	query := `
		DELETE FROM covers
		WHERE dir_id = :dirID
	`
	args := map[string]interface{}{
		"dirID": dirID,
	}

	_, err = tx.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirID).Msg("Failed to execute query to delete covers by directory")
		return err
	}

	log.Debug().Int("dirId", dirID).Msg("Covers deleted successfully by directory")
	return nil
}
