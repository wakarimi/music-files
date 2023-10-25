package song_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (r Repository) Create(tx *sqlx.Tx, song models.Song) (songId int, err error) {
	query := `
		INSERT INTO songs(dir_id, filename, extension, size_byte, duration_ms, bitrate_kbps, sample_rate_hz, channels_n, sha_256, last_content_update)
		VALUES (:dir_id, :filename, :extension, :size_byte, :duration_ms, :bitrate_kbps, :sample_rate_hz, :channels_n, :sha_256, CURRENT_TIMESTAMP)
		RETURNING song_id
	`
	rows, err := tx.NamedQuery(query, song)
	if err != nil {
		return 0, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err)
		}
	}(rows)

	if rows.Next() {
		if err := rows.Scan(&songId); err != nil {
			return 0, err
		}
	} else {
		err := fmt.Errorf("no id returned after song insert")
		return 0, err
	}

	return songId, nil
}
