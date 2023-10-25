package cover_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (r Repository) ReadAllByDir(tx *sqlx.Tx, dirId int) (covers []models.Cover, err error) {
	query := `
		SELECT * 
		FROM covers
		WHERE dir_id = :dir_id
	`
	args := map[string]interface{}{
		"dir_id": dirId,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err)
		}
	}(rows)

	for rows.Next() {
		var song models.Cover
		if err = rows.StructScan(&song); err != nil {
			log.Error().Err(err)
			return nil, err
		}
		covers = append(covers, song)
	}

	return covers, nil
}
