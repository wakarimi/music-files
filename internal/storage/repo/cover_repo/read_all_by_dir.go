package cover_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/cover"
)

func (r Repository) ReadAllByDir(tx *sqlx.Tx, dirID int) (covers []cover.Cover, err error) {
	log.Debug().Int("dirId", dirID).Msg("Reading covers by directory from database")

	query := `
		SELECT * 
		FROM covers
		WHERE dir_id = :dir_id
	`
	args := map[string]interface{}{
		"dir_id": dirID,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirID).Str("query", query).Msg("Failed to execute query to read cover by dirId")
		return nil, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close rows")
		}
	}(rows)

	for rows.Next() {
		var audioFile cover.Cover
		if err = rows.StructScan(&audioFile); err != nil {
			log.Error().Err(err).Int("dirId", dirID).Msg("Failed to get read result")
			return nil, err
		}
		covers = append(covers, audioFile)
	}

	log.Debug().Int("dirId", dirID).Int("countOfCoversInDir", len(covers)).Msg("Covers by dirId read successfully")
	return covers, nil
}
