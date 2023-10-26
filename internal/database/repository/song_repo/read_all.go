package song_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (r *Repository) ReadAll(tx *sqlx.Tx) (songs []models.Song, err error) {
	log.Debug().Msg("Fetching all songs")

	query := `
		SELECT * 
		FROM songs
	`
	err = tx.Select(&songs, query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch songs")
		return nil, err
	}

	log.Debug().Int("songsCount", len(songs)).Msg("All songs fetched successfully")
	return songs, nil
}
