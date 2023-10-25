package cover_repo

import (
	"github.com/jmoiron/sqlx"
)

func (r Repository) Delete(tx *sqlx.Tx, coverId int) (err error) {
	query := `
		DELETE FROM covers
		WHERE cover_id = :cover_id
	`
	args := map[string]interface{}{
		"cover_id": coverId,
	}
	_, err = tx.NamedExec(query, args)
	if err != nil {
		return err
	}

	return nil
}
