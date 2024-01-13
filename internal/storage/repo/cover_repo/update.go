package cover_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/cover"
)

func (r Repository) Update(tx *sqlx.Tx, coverId int, cover cover.Cover) (err error) {
	log.Debug().Int("coverId", coverId).Interface("cover", cover).Msg("Updating cover")

	query := `
		UPDATE covers
		SET dir_id = :dir_id, filename = :filename, extension = :extension, size_byte = :size_byte,
		    width_px = :width_px, height_px = :height_px, sha_256 = :sha_256, last_content_update = CURRENT_TIMESTAMP
		WHERE id = :id
	`

	cover.ID = coverId
	_, err = tx.NamedExec(query, cover)

	if err != nil {
		log.Error().Err(err).Int("coverId", coverId).Str("query", query).Msg("Failed to execute query to update cover")
		return err
	}

	log.Debug().Int("coverId", coverId).Msg("Cover updated successfully")
	return nil
}
