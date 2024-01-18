package cover_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/cover"
)

func (r Repository) Read(tx *sqlx.Tx, coverID int) (coverFile cover.Cover, err error) {
	log.Debug().Int("coverId", coverID).Msg("Reading cover file from database")

	query := `
		SELECT *
		FROM covers
		WHERE id = :id
	`
	args := map[string]interface{}{
		"id": coverID,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("coverId", coverID).Str("query", query).Msg("Failed to execute query to read cover file")
		return cover.Cover{}, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close rows")
		}
	}(rows)
	if rows.Next() {
		if err = rows.StructScan(&coverFile); err != nil {
			log.Error().Err(err).Int("coverId", coverID).Msg("Failed to get read result")
			return cover.Cover{}, err
		}
	} else {
		err := fmt.Errorf("No cover file found with id: %d", coverID)
		log.Error().Err(err).Int("coverId", coverID).Msg("Cover file not found")
		return cover.Cover{}, err
	}

	log.Debug().Interface("coverId", coverID).Msg("Cover file read successfully")
	return coverFile, nil
}
