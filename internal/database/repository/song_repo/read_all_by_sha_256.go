package song_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (r Repository) ReadAllBySha256(tx *sqlx.Tx, sha256 string) (songs []models.Song, err error) {
	query := `
		SELECT * 
		FROM songs
		WHERE sha_256 = :sha_256
	`
	args := map[string]interface{}{
		"sha_256": sha256,
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
