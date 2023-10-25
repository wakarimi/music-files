package song_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (r Repository) Read(tx *sqlx.Tx, songId int) (song models.Song, err error) {
	query := `
		SELECT *
		FROM songs
		WHERE song_id = :song_id
	`
	args := map[string]interface{}{
		"song_id": songId,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		return models.Song{}, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err)
		}
	}(rows)
	if rows.Next() {
		if err = rows.StructScan(&song); err != nil {
			return models.Song{}, err
		}
	} else {
		err := fmt.Errorf("No song found with songId: %d", songId)
		return models.Song{}, err
	}

	return song, nil
}
