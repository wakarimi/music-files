package cover_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) Delete(tx *sqlx.Tx, coverId int) (err error) {
	log.Debug().Int("coverId", coverId).Msg("Deleting cover from database")

	query := `
		DELETE FROM covers
		WHERE cover_id = :cover_id
	`
	args := map[string]interface{}{
		"cover_id": coverId,
	}
	_, err = tx.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("coverId", coverId).Msg("Failed to delete cover from database")
		return err
	}

	log.Debug().Int("coverId", coverId).Msg("Cover deleted from database successfully")
	return nil
}
