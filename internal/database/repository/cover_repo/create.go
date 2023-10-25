package cover_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (r Repository) Create(tx *sqlx.Tx, cover models.Cover) (coverId int, err error) {
	query := `
		INSERT INTO covers(dir_id, filename, extension, size_byte, width_px, height_px, sha_256, last_content_update)
		VALUES (:dir_id, :filename, :extension, :size_byte, :width_px, :height_px, :sha_256, CURRENT_TIMESTAMP)
		RETURNING cover_id
	`
	rows, err := tx.NamedQuery(query, cover)
	if err != nil {
		return 0, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err)
		}
	}(rows)

	if rows.Next() {
		if err := rows.Scan(&coverId); err != nil {
			return 0, err
		}
	} else {
		err := fmt.Errorf("no id returned after cover insert")
		return 0, err
	}

	return coverId, nil
}
