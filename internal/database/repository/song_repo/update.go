package song_repo

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/models"
)

func (r Repository) Update(tx *sqlx.Tx, songId int, song models.Song) (err error) {
	query := `
		UPDATE songs
		SET dir_id = :dir_id, filename = :filename, extension = :extension, size_byte = :size_byte,
		    duration_ms = :duration_ms, bitrate_kbps = :bitrate_kbps, sample_rate_hz = :sample_rate_hz,
		    channels_n = :channels_n, sha_256 = :sha_256, last_content_update = CURRENT_TIMESTAMP
		WHERE song_id = :song_id
	`

	song.SongId = songId
	_, err = tx.NamedExec(query, song)

	if err != nil {
		return err
	}

	return nil
}
