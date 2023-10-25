package cover_repo

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/models"
)

func (r Repository) Update(tx *sqlx.Tx, coverId int, cover models.Cover) (err error) {
	query := `
		UPDATE covers
		SET dir_id = :dir_id, filename = :filename, extension = :extension, size_byte = :size_byte,
		    width_px = :width_px, height_px = :height_px, sha_256 = :sha_256, last_content_update = CURRENT_TIMESTAMP
		WHERE cover_id = :cover_id
	`

	cover.CoverId = coverId
	_, err = tx.NamedExec(query, cover)

	if err != nil {
		return err
	}

	return nil
}
