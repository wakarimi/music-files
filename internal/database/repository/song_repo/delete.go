package song_repo

import (
	"github.com/jmoiron/sqlx"
)

func (r Repository) Delete(tx *sqlx.Tx, songId int) (err error) {
	query := `
		DELETE FROM songs
		WHERE song_id = :song_id
	`
	args := map[string]interface{}{
		"song_id": songId,
	}
	_, err = tx.NamedExec(query, args)
	if err != nil {
		return err
	}

	return nil
}
