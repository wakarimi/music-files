package cover_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model"
)

func (r Repository) ReadByDirAndName(tx *sqlx.Tx, dirId int, name string) (cover model.Cover, err error) {
	log.Debug().Int("dirId", dirId).Str("name", name).Msg("Reading cover from database")

	query := `
		SELECT *
		FROM covers
		WHERE dir_id = :dir_id
			AND filename = :name
	`
	args := map[string]interface{}{
		"dir_id": dirId,
		"name":   name,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Str("name", name).Str("query", query).Msg("Failed to execute query to read cover")
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
			log.Error().Err(err).Int("dirId", dirId).Str("name", name).Msg("Failed to get read result")
			return model.Cover{}, err
		}
	} else {
		err := fmt.Errorf("no cover found with dir_id: %d and name: %s", dirId, name)
		log.Error().Err(err).Int("dirId", dirId).Str("name", name).Msg("Cover not found")
		return model.Cover{}, err
	}

	log.Debug().Interface("cover", cover).Msg("Cover read successfully")
	return cover, nil
}
