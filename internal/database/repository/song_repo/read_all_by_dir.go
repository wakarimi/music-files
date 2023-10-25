package song_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (r Repository) ReadAllByDir(tx *sqlx.Tx, dirId int) (songs []models.Song, err error) {
	query := `
		SELECT * 
		FROM songs
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
		var song models.Song
		if err = rows.StructScan(&song); err != nil {
			log.Error().Err(err)
			return nil, err
		}
		songs = append(songs, song)
	}

	return songs, nil
}
