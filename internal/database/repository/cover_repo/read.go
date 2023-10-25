package cover_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (r Repository) Read(tx *sqlx.Tx, coverId int) (cover models.Cover, err error) {
	query := `
		SELECT *
		FROM covers
		WHERE cover_id = :cover_id
	`
	args := map[string]interface{}{
		"cover_id": coverId,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		return models.Cover{}, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err)
		}
	}(rows)
	if rows.Next() {
		if err = rows.StructScan(&cover); err != nil {
			return models.Cover{}, err
		}
	} else {
		err := fmt.Errorf("No cover found with coverId: %d", coverId)
		return models.Cover{}, err
	}

	return cover, nil
}
