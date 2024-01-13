package cover_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/cover"
)

func (r Repository) ReadByDirAndName(tx *sqlx.Tx, dirID int, name string) (coverFile cover.Cover, err error) {
	log.Debug().Int("dirId", dirID).Str("name", name).Msg("Reading cover from database")

	query := `
		SELECT *
		FROM covers
		WHERE dir_id = :dir_id
			AND filename = :name
	`
	args := map[string]interface{}{
		"dir_id": dirID,
		"name":   name,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirID).Str("name", name).Str("query", query).Msg("Failed to execute query to read cover")
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
			log.Error().Err(err).Int("dirId", dirID).Str("name", name).Msg("Failed to get read result")
			return cover.Cover{}, err
		}
	} else {
		err := fmt.Errorf("no cover found with dir_id: %d and name: %s", dirID, name)
		log.Error().Err(err).Int("dirId", dirID).Str("name", name).Msg("Cover not found")
		return cover.Cover{}, err
	}

	log.Debug().Interface("cover", coverFile).Msg("Cover read successfully")
	return coverFile, nil
}
