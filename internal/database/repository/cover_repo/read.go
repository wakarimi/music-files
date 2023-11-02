package cover_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model"
)

func (r Repository) Read(tx *sqlx.Tx, coverId int) (cover model.Cover, err error) {
	log.Debug().Int("coverId", coverId).Msg("Reading cover from database")

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
		log.Error().Err(err).Int("coverId", coverId).Str("query", query).Msg("Failed to execute query to read cover")
		return model.Cover{}, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close rows")
		}
	}(rows)
	if rows.Next() {
		if err = rows.StructScan(&cover); err != nil {
			log.Error().Err(err).Int("coverId", coverId).Msg("Failed to get read result")
			return model.Cover{}, err
		}
	} else {
		err := fmt.Errorf("no cover found with cover_id: %d", coverId)
		log.Error().Err(err).Int("coverId", coverId).Msg("Cover not found")
		return model.Cover{}, err
	}

	log.Debug().Interface("cover", cover).Msg("Cover read successfully")
	return cover, nil
}
