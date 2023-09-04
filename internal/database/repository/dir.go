package repository

import (
	"music-files/internal/database"
)

func InsertDir(path string) (dirId int, err error) {
	query := `
		INSERT INTO directories(path)
		VALUES ($1)
		RETURNING dir_id
	`
	err = database.Db.QueryRow(query, path).Scan(&dirId)
	if err != nil {
		return 0, err
	}

	return dirId, nil
}

func DirExist(path string) (exists bool, err error) {
	var count int

	query := `
		SELECT COUNT(dir_id)
		FROM directories
		WHERE path = $1
	`
	err = database.Db.Get(&count, query, path)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
