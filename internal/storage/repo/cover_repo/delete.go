package cover_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) Delete(tx *sqlx.Tx, coverID int) (err error) {
	log.Debug().Int("coverId", coverID).Msg("Deleting cover from database")

	query := `
		DELETE FROM covers
		WHERE id = :id
	`
	args := map[string]interface{}{
		"id": coverID,
	}
	_, err = tx.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("coverId", coverID).Msg("Failed to delete cover from database")
		return err
	}

	log.Debug().Int("coverId", coverID).Msg("Cover deleted from database successfully")
	return nil
}
