package cover_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/cover"
)

func (r Repository) Create(tx *sqlx.Tx, coverToCreate cover.Cover) (coverId int, err error) {
	log.Debug().Interface("cover", coverToCreate).Msg("Creating new cover in database")

	query := `
		INSERT INTO covers(dir_id, filename, extension, size_byte, width_px, height_px, sha_256, last_content_update)
		VALUES (:dir_id, :filename, :extension, :size_byte, :width_px, :height_px, :sha_256, CURRENT_TIMESTAMP)
		RETURNING id
	`
	rows, err := tx.NamedQuery(query, coverToCreate)
	if err != nil {
		log.Error().Err(err).Interface("cover", coverToCreate).Str("query", query).Msg("Failed to create cover in database")
		return 0, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close rows")
		}
	}(rows)

	if rows.Next() {
		if err := rows.Scan(&coverId); err != nil {
			log.Error().Err(err).Msg("Failed to scan coverId of created cover")
			return 0, err
		}
	} else {
		err := fmt.Errorf("no id returned after cover insert")
		log.Error().Err(err).Interface("cover", coverToCreate).Msg("No id returned after cover insert")
		return 0, err
	}

	log.Debug().Int("coverId", coverId).Interface("cover", coverToCreate).Msg("New cover in database created successfully")
	return coverId, nil
}
