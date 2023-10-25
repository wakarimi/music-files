package cover_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (r Repository) ReadByDirAndName(tx *sqlx.Tx, dirId int, name string) (song models.Cover, err error) {
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
		return models.Cover{}, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err)
		}
	}(rows)
	if rows.Next() {
		if err = rows.StructScan(&song); err != nil {
			return models.Cover{}, err
		}
	} else {
		err := fmt.Errorf("No cover found with dir_id: %d and name: %s", dirId, name)
		return models.Cover{}, err
	}

	return song, nil
}
